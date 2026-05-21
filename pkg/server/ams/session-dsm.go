// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"fmt"
	"session-dsm-grpc-plugin/pkg/client/ams"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
	"strings"

	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclient/game_session"
	"github.com/AccelByte/accelbyte-go-sdk/session-sdk/pkg/sessionclientmodels"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	ServerStatusReady = "READY"
	DSStatusAvailable = "AVAILABLE"
	AMSSource         = "AMS"
)

// AMSClientInterface defines the interface for AMS client operations.
type AMSClientInterface interface {
	ClaimDS(scope *envelope.Scope, params ams.DsClaimParameters) (ams.GameServer, error)
}

// SessionClientInterface defines the interface for session client operations.
type SessionClientInterface interface {
	AdminUpdateDSInformationShort(input *game_session.AdminUpdateDSInformationParams) error
}

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	AMSClient     AMSClientInterface
	SessionClient SessionClientInterface
}

// CreateGameSession implements synchronous game session creation using AMS
func (s *SessionDSM) CreateGameSession(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSession", "")
	defer scope.Finish()

	scope.Log.Info("Creating game session via AMS",
		"sessionID", req.SessionId,
		"namespace", req.Namespace,
		"deployment", req.Deployment,
		"regions", req.RequestedRegion,
	)

	// Prepare claim parameters based on request
	claimParams := ams.DsClaimParameters{
		SessionID: req.SessionId,
		Namespace: getGameNamespace(req.Namespace),
		ClaimKeys: strings.Split(req.Deployment, ","),
		Regions:   req.RequestedRegion,
	}

	// Validate parameters
	if err := claimParams.Validate(); err != nil {
		scope.TraceError(err)
		scope.Log.Error("Invalid claim parameters", "error", err)
		return nil, fmt.Errorf("invalid claim parameters: %w", err)
	}

	// Claim DS from AMS
	gameServer, err := s.AMSClient.ClaimDS(scope, claimParams)
	if err != nil {
		scope.TraceError(err)
		scope.Log.Error("Failed to claim DS from AMS", "error", err)
		return nil, fmt.Errorf("failed to claim DS from AMS: %w", err)
	}

	scope.Log.Info("Successfully claimed DS from AMS",
		"serverID", gameServer.PodName,
		"ip", gameServer.IP,
		"port", gameServer.Port,
		"region", gameServer.Region,
	)

	// Build response
	response := &sessiondsm.ResponseCreateGameSession{
		SessionId:     req.SessionId,
		Namespace:     req.Namespace,
		SessionData:   req.SessionData,
		Status:        ServerStatusReady,
		Ip:            gameServer.IP,
		Port:          int64(gameServer.Port),
		ServerId:      gameServer.PodName,
		Source:        AMSSource,
		Deployment:    gameServer.Deployment,
		Region:        gameServer.Region,
		ClientVersion: req.ClientVersion,
		GameMode:      req.GameMode,
		CreatedRegion: gameServer.Region,
	}

	return response, nil
}

// TerminateGameSession implements game session termination
// Note: AMS doesn't have a direct terminate API - servers are returned to the pool automatically
func (s *SessionDSM) TerminateGameSession(ctx context.Context, req *sessiondsm.RequestTerminateGameSession) (*sessiondsm.ResponseTerminateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "TerminateGameSession", "")
	defer scope.Finish()

	scope.Log.Info("Terminating game session",
		"sessionID", req.SessionId,
		"namespace", req.Namespace,
	)

	// For AMS, we don't need to explicitly terminate the server
	// The server will be automatically returned to the fleet pool
	// when the session ends in the session service

	response := &sessiondsm.ResponseTerminateGameSession{
		SessionId: req.SessionId,
		Namespace: req.Namespace,
		Success:   true,
		Reason:    "Session terminated successfully",
	}

	scope.Log.Info("Game session terminated", "sessionID", req.SessionId)

	return response, nil
}

// CreateGameSessionAsync implements asynchronous game session creation using AMS
func (s *SessionDSM) CreateGameSessionAsync(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSessionAsync, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSessionAsync", "")
	defer scope.Finish()

	scope.Log.Info("Creating game session asynchronously via AMS",
		"sessionID", req.SessionId,
		"namespace", req.Namespace,
		"deployment", req.Deployment,
		"regions", req.RequestedRegion,
	)

	// Return success immediately and process asynchronously
	response := &sessiondsm.ResponseCreateGameSessionAsync{
		Success: true,
		Message: "Game session creation initiated",
	}

	// Capture the trace context from the request span before launching the goroutine.
	// ContextWithRemoteSpanContext propagates the parent trace link without inheriting
	// the gRPC request context's cancellation, so the async work can outlive the RPC.
	parentCtx := oteltrace.ContextWithRemoteSpanContext(context.Background(), oteltrace.SpanContextFromContext(scope.Ctx))

	// Process DS claim asynchronously
	go func(requestData *sessiondsm.RequestCreateGameSession) {
		asyncScope := envelope.NewRootScope(parentCtx, "CreateGameSessionAsync.Worker", "")
		defer asyncScope.Finish()

		asyncScope.Log.Info("Starting async DS claim", "sessionID", requestData.SessionId)

		// Prepare claim parameters
		claimParams := ams.DsClaimParameters{
			SessionID: requestData.SessionId,
			Namespace: getGameNamespace(requestData.Namespace),
			ClaimKeys: strings.Split(requestData.Deployment, ","),
			Regions:   requestData.RequestedRegion,
		}

		// Validate parameters
		if err := claimParams.Validate(); err != nil {
			asyncScope.TraceError(err)
			asyncScope.Log.Error("Invalid claim parameters", "error", err)
			s.updateDSInformationWithError(asyncScope, requestData.Namespace, requestData.SessionId,
				fmt.Sprintf("Invalid claim parameters: %v", err))
			return
		}

		// Claim DS from AMS
		gameServer, err := s.AMSClient.ClaimDS(asyncScope, claimParams)
		if err != nil {
			asyncScope.TraceError(err)
			asyncScope.Log.Error("Failed to claim DS from AMS", "error", err)
			s.updateDSInformationWithError(asyncScope, requestData.Namespace, requestData.SessionId,
				fmt.Sprintf("Failed to claim DS: %v", err))
			return
		}

		asyncScope.Log.Info("Successfully claimed DS from AMS",
			"serverID", gameServer.PodName,
			"ip", gameServer.IP,
			"port", gameServer.Port,
			"region", gameServer.Region,
		)

		// Update session service with DS information
		s.updateDSInformation(asyncScope, requestData.Namespace, requestData.SessionId, gameServer)
	}(req)

	return response, nil
}

// updateDSInformation updates the session service with DS details
func (s *SessionDSM) updateDSInformation(scope *envelope.Scope, namespace, sessionID string, gameServer ams.GameServer) {
	port := int32(gameServer.Port)
	description := fmt.Sprintf("AMS Fleet: %s, Region: %s", gameServer.Deployment, gameServer.Region)

	dsStatusAvailable := DSStatusAvailable
	amsSource := AMSSource
	err := s.SessionClient.AdminUpdateDSInformationShort(&game_session.AdminUpdateDSInformationParams{
		Namespace: namespace,
		SessionID: sessionID,
		Context:   scope.Ctx,
		Body: &sessionclientmodels.ApimodelsUpdateGamesessionDSInformationRequest{
			Status:        &dsStatusAvailable,
			Source:        &amsSource,
			Port:          &port,
			ServerID:      &gameServer.PodName,
			IP:            &gameServer.IP,
			Deployment:    &gameServer.Deployment,
			Region:        &gameServer.Region,
			Description:   &description,
			CreatedRegion: &gameServer.Region,
		},
	})

	if err != nil {
		scope.Log.Error("Failed to update DS information", "error", err, "sessionID", sessionID)
		return
	}

	scope.Log.Info("Successfully updated DS information", "sessionID", sessionID, "serverID", gameServer.PodName)
}

// updateDSInformationWithError updates the session service with error status
func (s *SessionDSM) updateDSInformationWithError(scope *envelope.Scope, namespace, sessionID, errorMessage string) {
	status := "FAILED_TO_REQUEST"
	description := fmt.Sprintf("AMS claim failed: %s", errorMessage)
	amsSource := AMSSource
	err := s.SessionClient.AdminUpdateDSInformationShort(&game_session.AdminUpdateDSInformationParams{
		Namespace: namespace,
		SessionID: sessionID,
		Context:   scope.Ctx,
		Body: &sessionclientmodels.ApimodelsUpdateGamesessionDSInformationRequest{
			Status:      &status,
			Source:      &amsSource,
			Description: &description,
		},
	})

	if err != nil {
		scope.Log.Error("Failed to update DS error information", "error", err, "sessionID", sessionID)
	}
}

func getGameNamespace(namespace string) string {
	splittedString := strings.Split(namespace, "-")
	if len(splittedString) > 1 {
		return splittedString[1]
	}

	return namespace
}
