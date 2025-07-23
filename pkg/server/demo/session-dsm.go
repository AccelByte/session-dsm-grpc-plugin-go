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
	"time"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/session"
	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclient/game_session"
	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclientmodels"
)

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	Demo          *client.Client
	SessionClient *session.GameSessionService
	IamClient     *iam.OAuth20Service
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

		err := s.SessionClient.AdminUpdateDSInformationShort(&game_session.AdminUpdateDSInformationParams{
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
			scopeNew.Log.Error(err)
		}
	}(req)

	return responses, nil
}
