// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"errors"
	"session-dsm-grpc-plugin/pkg/client/gcpvm"
	"session-dsm-grpc-plugin/pkg/constants"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/session"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"time"
)

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	ClientGCP *gcpvm.GCPVM
}

func (s *SessionDSM) CreateGameSession(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSession", "")
	defer scope.Finish()
	var response *sessiondsm.ResponseCreateGameSession
	var err error

	if len(req.RequestedRegion) == 0 {
		return nil, errors.New("need provide requested region")
	}

	for _, region := range req.RequestedRegion {
		response, err = s.ClientGCP.CreateGameSession(scope, "", req.SessionId, req.SessionData, region, int(req.MaximumPlayer), req.Namespace, req.Deployment)
		if err != nil {
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	responses := &sessiondsm.ResponseCreateGameSession{
		SessionId:     req.SessionId,
		Namespace:     req.Namespace,
		Deployment:    req.Deployment,
		SessionData:   req.SessionData,
		Status:        constants.ServerStatusReady,
		Ip:            response.Ip,
		Port:          int64(response.Port),
		ServerId:      response.ServerId,
		Source:        constants.GameServerSourceGCP,
		Region:        response.Region,
		ClientVersion: req.ClientVersion,
		GameMode:      req.GameMode,
		CreatedRegion: response.CreatedRegion,
	}
	return responses, err
}

func (s *SessionDSM) TerminateGameSession(ctx context.Context, req *sessiondsm.RequestTerminateGameSession) (*sessiondsm.ResponseTerminateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "TerminateGameSession", "")
	defer scope.Finish()
	response, err := s.ClientGCP.TerminateGameSession(scope, req.SessionId, req.Namespace, req.Zone)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (s *SessionDSM) CreateGameSessionAsync(ctx context.Context,
	req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSessionAsync, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSessionAsync", "")
	defer scope.Finish()

	responses := &sessiondsm.ResponseCreateGameSessionAsync{
		Success: true,
		Message: "success",
	}

	// example for put DS Information after 2 second
	go func(data *sessiondsm.RequestCreateGameSession) {
		scopeNew := envelope.NewRootScope(context.Background(), "SendData", "")
		defer scopeNew.Finish()
		time.Sleep(2 * time.Second)
		_, err := s.ClientGCP.UpdateDSInformation(scopeNew,
			&session.UpdateGamesessionDSInformationRequest{
				Status:      session.DSStatusAvailable,
				Port:        1223,
				ServerID:    "123455",
				IP:          "192.168.1.1",
				Description: "testing",
			}, data.Namespace, data.SessionId,
		)
		if err != nil {
			scopeNew.Log.Error(err)
		}
	}(req)
	return responses, nil
}
