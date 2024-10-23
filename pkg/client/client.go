// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package client

import (
	"session-dsm-grpc-plugin/pkg/client/model"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

type Client interface {
	CreateGameSession(scope *envelope.Scope, fleetAlias, sessionID, sessionData, location string, maxPlayer int) (*model.ResponseCreateGameSession, error)
	TerminateGameSession(scope *envelope.Scope, SessionID, Namespace string) (*model.ResponseTerminateGameSession, error)
}
