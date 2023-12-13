package server

import (
	"context"
	"session-dsm-grpc-plugin/pkg/awsgamelift"
	"session-dsm-grpc-plugin/pkg/constants"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"session-dsm-grpc-plugin/pkg/utils/envelope"
)

type SessionDSM struct {
	sessiondsm.UnimplementedSessionDsmServer
	ClientGamelift *awsgamelift.AwsGamelift
}

func (s *SessionDSM) CreateGameSession(ctx context.Context, req *sessiondsm.RequestCreateGameSession) (*sessiondsm.ResponseCreateGameSession, error) {
	scope := envelope.NewRootScope(ctx, "CreateGameSession", "")
	defer scope.Finish()
	response, err := s.ClientGamelift.CreateGameSession(scope, req.FleetAlias, req.SessionId, req.SessionData, req.RequestedRegion, int(req.MaximumPlayer))
	if err != nil {
		return nil, err
	}

	responses := &sessiondsm.ResponseCreateGameSession{
		SessionId:   req.SessionId,
		Namespace:   req.Namespace,
		FleetAlias:  req.FleetAlias,
		SessionData: req.SessionData,
		Status:      constants.ServerStatusReady,
		Ip:          response.IPAddress,
		Port:        int64(response.Port),
		ServerId:    response.GameSessionARN,
		Source:      constants.GameServerSourceGamelift,
		Deployment:  response.FleetID,
		Region:      response.Location,
	}
	return responses, err
}
