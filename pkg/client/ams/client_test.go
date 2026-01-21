// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package ams

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDsClaimParameters_Validate(t *testing.T) {
	tests := []struct {
		name      string
		params    DsClaimParameters
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid parameters",
			params: DsClaimParameters{
				SessionID: "session-123",
				Namespace: "mygame",
				ClaimKeys: []string{"production-fleet"},
				Regions:   []string{"us-west-2"},
			},
			wantError: false,
		},
		{
			name: "missing namespace",
			params: DsClaimParameters{
				SessionID: "session-123",
				Namespace: "",
				ClaimKeys: []string{"production-fleet"},
				Regions:   []string{"us-west-2"},
			},
			wantError: true,
			errorMsg:  "namespace is required",
		},
		{
			name: "missing claim keys",
			params: DsClaimParameters{
				SessionID: "session-123",
				Namespace: "mygame",
				ClaimKeys: []string{},
				Regions:   []string{"us-west-2"},
			},
			wantError: true,
			errorMsg:  "at least one claim key is required",
		},
		{
			name: "missing regions",
			params: DsClaimParameters{
				SessionID: "session-123",
				Namespace: "mygame",
				ClaimKeys: []string{"production-fleet"},
				Regions:   []string{},
			},
			wantError: true,
			errorMsg:  "at least one region is required",
		},
		{
			name: "all fields missing",
			params: DsClaimParameters{
				SessionID: "session-123",
				Namespace: "",
				ClaimKeys: []string{},
				Regions:   []string{},
			},
			wantError: true,
			errorMsg:  "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if tt.wantError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_ClaimDS_Success(t *testing.T) {
	// Create mock AMS server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Contains(t, r.URL.Path, "/v1/namespaces/mygame/servers/claim")

		// Verify headers
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Verify request body
		var reqBody claimByKeysReq
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)
		assert.Equal(t, []string{"production-fleet"}, reqBody.ClaimKeys)
		assert.Equal(t, []string{"us-west-2"}, reqBody.Regions)
		assert.Equal(t, "session-123", reqBody.SessionID)

		// Send success response
		response := claimResponse{
			ServerID: "server-abc123",
			IP:       "192.168.1.100",
			Ports: map[string]int{
				"default": 7777,
				"query":   7778,
			},
			Region:  "us-west-2",
			FleetID: "fleet-xyz",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create client
	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	// Create scope
	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	// Test parameters
	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "mygame",
		ClaimKeys: []string{"production-fleet"},
		Regions:   []string{"us-west-2"},
	}

	// Execute
	gameServer, err := client.ClaimDS(scope, params)

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, "server-abc123", gameServer.PodName)
	assert.Equal(t, "192.168.1.100", gameServer.IP)
	assert.Equal(t, 7777, gameServer.Port)
	assert.Equal(t, "us-west-2", gameServer.Region)
	assert.Equal(t, "fleet-xyz", gameServer.FleetID)
	assert.Equal(t, "claimed", gameServer.StatusText)
	assert.Equal(t, "session-123", gameServer.SessionID)
	assert.Equal(t, "mygame", gameServer.Namespace)
	assert.NotNil(t, gameServer.Ports)
	assert.Equal(t, 7778, gameServer.Ports["query"])
}

func TestClient_ClaimDS_ValidationError(t *testing.T) {
	client := New(
		"http://localhost",
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	// Invalid parameters (missing namespace)
	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "",
		ClaimKeys: []string{"production-fleet"},
		Regions:   []string{"us-west-2"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")
	assert.Equal(t, "", gameServer.PodName)
}

func TestClient_ClaimDS_NotFound(t *testing.T) {
	// Create mock AMS server returning 404
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := map[string]string{
			"errorMessage": "no matching DS available to claim",
			"traceId":      "trace-456",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResp)
	}))
	defer mockServer.Close()

	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "mygame",
		ClaimKeys: []string{"production-fleet"},
		Regions:   []string{"us-west-2"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no matching DS available to claim")
	assert.Contains(t, err.Error(), "trace-456")
	assert.Equal(t, "", gameServer.PodName)
}

func TestClient_ClaimDS_BadRequest(t *testing.T) {
	// Create mock AMS server returning 400
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := map[string]string{
			"errorMessage": "invalid claim keys",
			"traceId":      "trace-789",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
	}))
	defer mockServer.Close()

	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "mygame",
		ClaimKeys: []string{"invalid-fleet"},
		Regions:   []string{"us-west-2"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "malformed request")
	assert.Contains(t, err.Error(), "invalid claim keys")
	assert.Equal(t, "", gameServer.PodName)
}

func TestClient_ClaimDS_ServerError(t *testing.T) {
	// Create mock AMS server returning 500
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := map[string]string{
			"errorMessage": "internal server error",
			"traceId":      "trace-500",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResp)
	}))
	defer mockServer.Close()

	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "mygame",
		ClaimKeys: []string{"production-fleet"},
		Regions:   []string{"us-west-2"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server error")
	assert.Contains(t, err.Error(), "internal server error")
	assert.Equal(t, "", gameServer.PodName)
}

func TestClient_ClaimDS_InvalidJSON(t *testing.T) {
	// Create mock AMS server returning invalid JSON
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json{"))
	}))
	defer mockServer.Close()

	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	params := DsClaimParameters{
		SessionID: "session-123",
		Namespace: "mygame",
		ClaimKeys: []string{"production-fleet"},
		Regions:   []string{"us-west-2"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to unmarshal")
	assert.Equal(t, "", gameServer.PodName)
}

func TestClient_ClaimDS_NoDefaultPort(t *testing.T) {
	// Create mock AMS server without default port
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := claimResponse{
			ServerID: "server-xyz",
			IP:       "10.0.0.1",
			Ports: map[string]int{
				"custom": 9999,
			},
			Region:  "us-east-1",
			FleetID: "fleet-abc",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	client := New(
		mockServer.URL,
		func() string { return "test-token" },
		10*time.Second,
	)

	scope := envelope.NewRootScope(context.Background(), "TestClaimDS", "test-flight-id")

	params := DsClaimParameters{
		SessionID: "session-456",
		Namespace: "testgame",
		ClaimKeys: []string{"test-fleet"},
		Regions:   []string{"us-east-1"},
	}

	gameServer, err := client.ClaimDS(scope, params)

	assert.NoError(t, err)
	assert.Equal(t, "server-xyz", gameServer.PodName)
	assert.Equal(t, 0, gameServer.Port) // No default port
	assert.Equal(t, 9999, gameServer.Ports["custom"])
}

func TestParseErrorResponse(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		expected string
	}{
		{
			name: "valid error response",
			body: []byte(`{"errorMessage":"test error","traceId":"trace-123"}`),
			expected: "; error message: test error; trace ID: trace-123",
		},
		{
			name:     "invalid json",
			body:     []byte(`invalid json`),
			expected: "invalid json",
		},
		{
			name:     "empty response",
			body:     []byte(`{}`),
			expected: "{}",
		},
		{
			name:     "only error message",
			body:     []byte(`{"errorMessage":"error only"}`),
			expected: "; error message: error only; trace ID: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseErrorResponse(tt.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNew(t *testing.T) {
	tokenFunc := func() string { return "my-token" }
	timeout := 30 * time.Second

	client := New("https://ams.example.com", tokenFunc, timeout)

	assert.NotNil(t, client)
	assert.Equal(t, "https://ams.example.com", client.baseURL)
	assert.Equal(t, timeout, client.requestTimeout)
	assert.NotNil(t, client.httpClient)
	assert.NotNil(t, client.getTokenFunc)
	assert.Equal(t, "my-token", client.getTokenFunc())
}
