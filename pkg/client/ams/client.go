// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package ams

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"time"

	"github.com/cenkalti/backoff"
)

const (
	maxBackoffCount = uint64(5)
)

// Client represents an AMS client for claiming dedicated servers
type Client struct {
	baseURL        string
	httpClient     *http.Client
	getTokenFunc   func() string
	requestTimeout time.Duration
}

// New creates a new AMS client
func New(baseURL string, getTokenFunc func() string, requestTimeout time.Duration) *Client {
	return &Client{
		baseURL:        baseURL,
		httpClient:     &http.Client{Timeout: requestTimeout},
		getTokenFunc:   getTokenFunc,
		requestTimeout: requestTimeout,
	}
}

// GameServer represents a claimed game server from AMS
type GameServer struct {
	SessionID  string
	Namespace  string
	PodName    string // ServerID
	IP         string
	Port       int
	Ports      map[string]int
	Region     string
	Deployment string
	FleetID    string
	StatusText string
}

// DsClaimParameters contains parameters for claiming a DS
type DsClaimParameters struct {
	SessionID string
	Namespace string
	ClaimKeys []string
	Regions   []string
}

// Validate validates the claim parameters
func (p DsClaimParameters) Validate() error {
	var errs []string

	if p.Namespace == "" {
		errs = append(errs, "namespace is required")
	}

	if len(p.ClaimKeys) == 0 {
		errs = append(errs, "at least one claim key is required")
	}

	if len(p.Regions) == 0 {
		errs = append(errs, "at least one region is required")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation error: %v", errs)
	}

	return nil
}

type claimByKeysReq struct {
	ClaimKeys []string `json:"claimKeys"`
	Regions   []string `json:"regions"`
	SessionID string   `json:"sessionId"`
}

type claimResponse struct {
	ServerID string         `json:"serverId"`
	IP       string         `json:"ip"`
	Ports    map[string]int `json:"ports"`
	Region   string         `json:"region"`
	FleetID  string         `json:"fleetId"`
}

// ClaimDS claims a dedicated server from AMS
func (c *Client) ClaimDS(scope *envelope.Scope, parameters DsClaimParameters) (GameServer, error) {
	scope.Log.Info("Claiming DS from AMS",
		"namespace", parameters.Namespace,
		"sessionID", parameters.SessionID,
		"regions", parameters.Regions,
		"claimKeys", parameters.ClaimKeys,
	)

	// Validate parameters
	if err := parameters.Validate(); err != nil {
		return GameServer{}, err
	}

	request := claimByKeysReq{
		ClaimKeys: parameters.ClaimKeys,
		Regions:   parameters.Regions,
		SessionID: parameters.SessionID,
	}

	scope.Log.Info("AMS ClaimDS request details",
		"baseURL", c.baseURL,
		"namespace", parameters.Namespace,
		"sessionID", parameters.SessionID,
		"claimKeys", parameters.ClaimKeys,
		"regions", parameters.Regions,
		"requestBody", request,
	)

	var (
		body []byte
		resp *http.Response
		err  error
	)

	// Retry logic for claiming DS
	err = backoff.RetryNotify(func() error {
		url := fmt.Sprintf("%s/ams/v1/namespaces/%s/servers/claim", c.baseURL, parameters.Namespace)
		scope.Log.Info("AMS ClaimDS sending PUT request", "url", url)
		resp, body, err = c.put(scope, url, request)
		if err != nil {
			scope.Log.Error("Failed to make claim request", "error", err)
			return backoff.Permanent(fmt.Errorf("unable to make request to claim DS: %w", err))
		}

		scope.Log.Info("AMS ClaimDS response received",
			"statusCode", resp.StatusCode,
			"responseBody", string(body),
		)

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			errMsg := fmt.Sprintf("client error claiming DS (status %d)", resp.StatusCode)
			if resp.StatusCode == http.StatusBadRequest {
				errMsg = "malformed request"
			} else if resp.StatusCode == http.StatusNotFound {
				errMsg = "no matching DS available to claim"
			}
			errMsg += parseErrorResponse(body)
			scope.Log.Error(errMsg, "statusCode", resp.StatusCode)
			return backoff.Permanent(errors.New(errMsg))
		}

		if resp.StatusCode >= 500 {
			errMsg := fmt.Sprintf("server error claiming DS (status %d)", resp.StatusCode)
			errMsg += parseErrorResponse(body)
			scope.Log.Error(errMsg, "statusCode", resp.StatusCode)
			return backoff.Permanent(errors.New(errMsg))
		}

		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxBackoffCount),
		func(_ error, _ time.Duration) {
			scope.Log.Info("Retrying claim DS from AMS", "sessionID", parameters.SessionID)
		})

	if err != nil {
		return GameServer{}, fmt.Errorf("unable to claim DS from AMS: %w", err)
	}

	// Parse response
	scope.Log.Info("AMS ClaimDS claim successful, parsing response", "rawBody", string(body))
	var result claimResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return GameServer{}, fmt.Errorf("unable to unmarshal claim response: %w", err)
	}

	scope.Log.Info("AMS ClaimDS parsed response",
		"serverID", result.ServerID,
		"ip", result.IP,
		"ports", result.Ports,
		"region", result.Region,
		"fleetID", result.FleetID,
	)

	// Get default port
	defaultPort := 0
	if port, ok := result.Ports["default"]; ok {
		defaultPort = port
	}
	scope.Log.Info("AMS ClaimDS resolved default port", "defaultPort", defaultPort, "allPorts", result.Ports)

	// Build GameServer response
	gameServer := GameServer{
		SessionID:  parameters.SessionID,
		Namespace:  parameters.Namespace,
		PodName:    result.ServerID,
		IP:         result.IP,
		Port:       defaultPort,
		Ports:      result.Ports,
		Region:     result.Region,
		FleetID:    result.FleetID,
		StatusText: "claimed",
		Deployment: "", // Will be filled from fleet details if needed
	}

	scope.Log.Info("Successfully claimed DS from AMS",
		"serverID", gameServer.PodName,
		"ip", gameServer.IP,
		"port", gameServer.Port,
		"region", gameServer.Region,
		"fleetID", gameServer.FleetID,
	)

	return gameServer, nil
}

// put makes a PUT request to AMS
func (c *Client) put(scope *envelope.Scope, url string, requestBody any) (*http.Response, []byte, error) {
	token := c.getTokenFunc()

	tokenPreview := "(empty)"
	if len(token) > 10 {
		tokenPreview = token[:10] + "..."
	}
	scope.Log.Info("AMS PUT request preparing",
		"url", url,
		"tokenPreview", tokenPreview,
		"timeout", c.requestTimeout,
	)

	ctx, cancel := context.WithTimeout(scope.Ctx, c.requestTimeout)
	defer cancel()

	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, err
	}

	scope.Log.Info("AMS PUT request body", "body", string(b))

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				scope.Log.Error("unable to close body: ", err.Error())
			}
		}(resp.Body)
	}

	if err != nil {
		return nil, nil, err
	}

	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBodyBytes, nil
}

// parseErrorResponse parses AMS error response
func parseErrorResponse(body []byte) string {
	type errorResponse struct {
		TraceID      string `json:"traceId"`
		ErrorMessage string `json:"errorMessage"`
	}

	var response errorResponse
	err := json.Unmarshal(body, &response)
	if err != nil || (response.TraceID == "" && response.ErrorMessage == "") {
		return string(body)
	}

	return fmt.Sprintf("; error message: %s; trace ID: %s", response.ErrorMessage, response.TraceID)
}
