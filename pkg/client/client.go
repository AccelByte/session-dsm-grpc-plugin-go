package client

import (
	"session-dsm-grpc-plugin/pkg/client/model"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

type Client interface {
	CreateGameSession(scope *envelope.Scope, fleetAlias, sessionID, sessionData, location string, maxPlayer int) (*model.ResponseCreateGameSession, error)
	TerminateGameSession(scope *envelope.Scope, SessionID, Namespace string) (*model.ResponseTerminateGameSession, error)
}
