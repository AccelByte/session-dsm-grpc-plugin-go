// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package model

type ResponseCreateGameSession struct {
	SessionId     string
	Namespace     string
	SessionData   string
	Status        string
	Ip            string
	Port          int64
	ServerId      string
	Source        string
	Deployment    string
	Region        string
	ClientVersion string
	GameMode      string
}

type ResponseTerminateGameSession struct {
	SessionId string
	Namespace string
	Success   bool
	Reason    string
}
