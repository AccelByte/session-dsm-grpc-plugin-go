// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package session

import "errors"

const (
	DSStatusAvailable       = "AVAILABLE"
	DSStatusFailedToRequest = "FAILED_TO_REQUEST"
)

type UpdateGamesessionDSInformationRequest struct {
	Status        string `json:"status"`
	Source        string `json:"source"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	ServerID      string `json:"serverId"`
	Deployment    string `json:"deployment"`
	Region        string `json:"region"`
	ClientVersion string `json:"clientVersion"`
	GameMode      string `json:"gameMode"`
	Description   string `json:"description"`
	CreatedRegion string `json:"createdRegion"`
}

func (u *UpdateGamesessionDSInformationRequest) Validate() error {
	if !(u.Status == DSStatusAvailable ||
		u.Status == DSStatusFailedToRequest) {
		return errors.New("Status DS is not Valid, the Valid Status is AVAILABLE or FAILED_TO_REQUEST")
	}

	if u.Port == 0 {
		return errors.New("Port is not Valid, the Port cannot be 0")
	}

	if u.ServerID == "" {
		return errors.New("invalid Server ID cannot be empty")
	}

	return nil
}
