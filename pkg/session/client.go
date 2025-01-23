// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package session

import (
	"errors"
	"fmt"
	"net/http"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
)

type SessionClient struct {
	iamClient      *iam.OAuth20Service
	defaultBaseURL string
	requestTimeout time.Duration
}

func New(iamClient *iam.OAuth20Service, baseURL string, reqTimeout time.Duration) *SessionClient {
	return &SessionClient{
		iamClient:      iamClient,
		defaultBaseURL: baseURL,
		requestTimeout: reqTimeout,
	}
}

//nolint:cyclop,funlen
func (s *SessionClient) RequestAdminUpdateDSInformation(scope *envelope.Scope,
	request *UpdateGamesessionDSInformationRequest, namespace, sessionID string) (int, error) {
	var body []byte
	var resp *http.Response
	var err error

	if err = request.Validate(); err != nil {
		return 0, err
	}

	errRequestUpdate := fmt.Sprintf("Namespace %s Session ID %s", namespace, sessionID)

	url := fmt.Sprintf("%s/session/v1/admin/namespaces/%s/gamesessions/%s/dsinformation", s.defaultBaseURL, namespace, sessionID)
	resp, body, err = s.put(scope, url, request) //nolint:bodyclose
	if err != nil {
		errString := "unable to make request to update DS Information for " + errRequestUpdate + ": " + err.Error()
		err := errors.New(errString)
		scope.TraceError(err)
		return 0, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		errString := "unable to update DS Information for " + errRequestUpdate + ": "
		switch resp.StatusCode {
		case http.StatusBadRequest:
			errString += resp.Status
		default:
			errString += "unexpected status: " + resp.Status
		}
		errString += parseErrorResponse(body)
		return resp.StatusCode, errors.New(errString)
	}

	if resp.StatusCode >= 500 {
		errString := "unable to update DS Information for " + errRequestUpdate + ": " + resp.Status
		errString += parseErrorResponse(body)
		err := errors.New(errString)
		scope.TraceError(err)
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}
