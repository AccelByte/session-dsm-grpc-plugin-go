// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package awsgamelift

type GameSessionResult struct {
	FleetID        string
	GameSessionARN string
	IPAddress      string
	Port           int
	Location       string
}