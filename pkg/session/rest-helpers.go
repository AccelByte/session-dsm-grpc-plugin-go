// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package session

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

//nolint:nonamedreturns
func (s *SessionClient) put(scope *envelope.Scope, url string, requestBody interface{}) (response *http.Response, body []byte, err error) {
	token, err := s.iamClient.GetToken()
	if err != nil {
		scope.Log.Error(err)
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(scope.Ctx, s.requestTimeout)
	defer cancel()

	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}
	h := http.Header{}
	h.Add("Authorization", "Bearer "+token)
	h.Add("Content-Type", "application/json")
	req.Header = h

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBodyBytes, nil
}
