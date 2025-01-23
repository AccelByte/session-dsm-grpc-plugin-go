// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package session

import (
	"encoding/json"
	"fmt"

	"github.com/AccelByte/go-restful-plugins/v3/pkg/response"
)

func parseErrorResponse(body []byte) string {
	var errResp response.Error
	err := json.Unmarshal(body, &errResp)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(": error message %s: error code %d", errResp.ErrorMessage, errResp.ErrorCode)
}
