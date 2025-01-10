// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"fmt"
	"session-dsm-grpc-plugin/pkg/client"
	"session-dsm-grpc-plugin/pkg/constants"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	Demo *client.Client
}

func (s *SessionDSM) CreateGameSession(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSession", "")
	defer scope.Finish()
	var err error

	responses := &sessiondsm.ResponseCreateGameSession{
		SessionId:     req.SessionId,
		Namespace:     req.Namespace,
		Deployment:    req.Deployment,
		SessionData:   req.SessionData,
		Status:        constants.ServerStatusReady,
		Ip:            "10.10.10.11",
		Port:          int64(8080),
		ServerId:      fmt.Sprintf("demo-local-%s", req.SessionId),
		Source:        "DEMO",
		Region:        req.RequestedRegion[0],
		ClientVersion: req.ClientVersion,
		GameMode:      req.GameMode,
		CreatedRegion: req.RequestedRegion[0],
	}
	return responses, err
}

func (s *SessionDSM) TerminateGameSession(ctx context.Context, req *sessiondsm.RequestTerminateGameSession) (*sessiondsm.ResponseTerminateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "TerminateGameSession", "")
	defer scope.Finish()
	var err error

	responses := &sessiondsm.ResponseTerminateGameSession{
		SessionId: req.SessionId,
		Namespace: req.Namespace,
	}
	return responses, err
}
