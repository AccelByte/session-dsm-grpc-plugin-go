// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package awsgamelift

import (
	"encoding/base64"
	"encoding/json"
)

func ConvertToJSONBase64(data any) (string, error) {
	jsonByte, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(jsonByte), nil
}
