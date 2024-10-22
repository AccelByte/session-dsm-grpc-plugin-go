// Copyright (c) 2023 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package awsgamelift

import (
	"errors"
	"session-dsm-grpc-plugin/pkg/utils/envelope"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

const (
	TypeCreateGameliftGamesession        = "gamelift_create_gamesession"
	TypeGameliftGamesessionCreated       = "gamelift_gamesession_created"
	TypeCustomeCreateGameliftGamesession = "gamelift_custome_create_gamesession"
)

type AwsGamelift struct {
	credential *credentials.Credentials
	region     string
}

func New(credential *credentials.Credentials, gameliftRegion string) *AwsGamelift {
	return &AwsGamelift{
		credential: credential,
		region:     gameliftRegion,
	}
}

func (a *AwsGamelift) CreateGameSession(rootScope *envelope.Scope, fleetAlias, sessionID, sessionData, location string, maxPlayer int) (*GameSessionResult, error) {
	scope := rootScope.NewChildScope("awsgamelift.CreateGameSession")
	defer scope.Finish()
	sess, err := session.NewSession(&aws.Config{
		Credentials: a.credential,
		Region:      &a.region,
	})
	if err != nil {
		return nil, err
	}
	srv := gamelift.New(sess, aws.NewConfig().WithRegion(a.region))
	createRequest := &gamelift.CreateGameSessionInput{
		AliasId:                   &fleetAlias,
		GameSessionData:           &sessionData,
		IdempotencyToken:          &sessionID,
		MaximumPlayerSessionCount: aws.Int64(int64(maxPlayer)),
	}
	if location != "" {
		createRequest.Location = &location
	}
	result, err := srv.CreateGameSessionWithContext(scope.Ctx, createRequest)
	if err != nil {
		return nil, err
	}
	if result.GameSession == nil {
		return nil, errors.New("returned nil gamesession")
	}
	return &GameSessionResult{
		FleetID:        *result.GameSession.FleetId,
		GameSessionARN: *result.GameSession.GameSessionId,
		IPAddress:      *result.GameSession.IpAddress,
		Port:           int(*result.GameSession.Port),
		Location:       *result.GameSession.Location,
	}, nil
}
