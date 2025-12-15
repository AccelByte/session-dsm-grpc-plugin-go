// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"errors"
	"session-dsm-grpc-plugin/pkg/client/awsgamelift"
	"session-dsm-grpc-plugin/pkg/constants"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclient/game_session"
	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclientmodels"
)

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	ClientGamelift *awsgamelift.AwsGamelift
}

func (s *SessionDSM) CreateGameSession(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSession", "")
	defer scope.Finish()
	var response *awsgamelift.GameSessionResult
	var err error

	if len(req.RequestedRegion) == 0 {
		return nil, errors.New("need provide requested region")
	}

	for _, region := range req.RequestedRegion {
		response, err = s.ClientGamelift.CreateGameSession(scope, req.Deployment, req.SessionId, req.SessionData, region, int(req.MaximumPlayer))
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
		Deployment:    response.FleetID,
		SessionData:   req.SessionData,
		Status:        constants.ServerStatusReady,
		Ip:            response.IPAddress,
		Port:          int64(response.Port),
		ServerId:      response.GameSessionARN,
		Source:        constants.GameServerSourceGamelift,
		Region:        response.Location,
		ClientVersion: req.ClientVersion,
		GameMode:      req.GameMode,
		CreatedRegion: response.Location,
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

		var port int32 = 1223
		serverID := "123455"
		ip := "192.168.1.1"
		description := "testing"

		err := s.ClientGamelift.UpdateDSInformation(scopeNew, &game_session.AdminUpdateDSInformationParams{
			Namespace: data.Namespace,
			SessionID: data.SessionId,
			Context:   scopeNew.Ctx,
			Body: &sessionclientmodels.ApimodelsUpdateGamesessionDSInformationRequest{
				Status:      &constants.DSStatusAvailable,
				Port:        &port,
				ServerID:    &serverID,
				IP:          &ip,
				Description: &description,
			},
		})
		if err != nil {
			scopeNew.Log.Error("Failed to update DS information", "error", err)
		}
	}(req)

	return responses, nil
}
