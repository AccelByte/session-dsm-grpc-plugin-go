// Copyright (c) 2025 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package server

import (
	"context"
	"errors"
	"session-dsm-grpc-plugin/pkg/client/ams"
	sessiondsm "session-dsm-grpc-plugin/pkg/pb"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAMSClient is a mock implementation of the AMS client
type MockAMSClient struct {
	mock.Mock
}

func (m *MockAMSClient) ClaimDS(scope interface{}, params ams.DsClaimParameters) (ams.GameServer, error) {
	args := m.Called(scope, params)
	return args.Get(0).(ams.GameServer), args.Error(1)
}

// MockSessionClient is a mock implementation of the session client
type MockSessionClient struct {
	mock.Mock
}

func (m *MockSessionClient) AdminUpdateDSInformationShort(params interface{}) error {
	args := m.Called(params)
	return args.Error(0)
}

func TestSessionDSM_CreateGameSession_Success(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	// Setup expected call
	expectedGameServer := ams.GameServer{
		SessionID:  "session-123",
		Namespace:  "mygame",
		PodName:    "server-abc123",
		IP:         "192.168.1.100",
		Port:       7777,
		Region:     "us-west-2",
		Deployment: "production-fleet",
		StatusText: "claimed",
	}

	mockAMSClient.On("ClaimDS", mock.Anything, mock.MatchedBy(func(params ams.DsClaimParameters) bool {
		return params.SessionID == "session-123" &&
			params.Namespace == "mygame" &&
			len(params.ClaimKeys) == 1 &&
			params.ClaimKeys[0] == "production-fleet" &&
			len(params.Regions) == 1 &&
			params.Regions[0] == "us-west-2"
	})).Return(expectedGameServer, nil)

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	// Create request
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-123",
		Namespace:       "mygame",
		Deployment:      "production-fleet",
		SessionData:     `{"mapName":"Desert"}`,
		RequestedRegion: []string{"us-west-2"},
		MaximumPlayer:   10,
		ClientVersion:   "1.0.0",
		GameMode:        "deathmatch",
		Secret:          "secret-key",
	}

	// Execute
	response, err := server.CreateGameSession(context.Background(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "session-123", response.SessionId)
	assert.Equal(t, "mygame", response.Namespace)
	assert.Equal(t, ServerStatusReady, response.Status)
	assert.Equal(t, "192.168.1.100", response.Ip)
	assert.Equal(t, int64(7777), response.Port)
	assert.Equal(t, "server-abc123", response.ServerId)
	assert.Equal(t, AMSSource, response.Source)
	assert.Equal(t, "production-fleet", response.Deployment)
	assert.Equal(t, "us-west-2", response.Region)
	assert.Equal(t, "us-west-2", response.CreatedRegion)
	assert.Equal(t, "1.0.0", response.ClientVersion)
	assert.Equal(t, "deathmatch", response.GameMode)

	mockAMSClient.AssertExpectations(t)
}

func TestSessionDSM_CreateGameSession_ClaimFailure(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	// Setup expected call to return error
	mockAMSClient.On("ClaimDS", mock.Anything, mock.Anything).
		Return(ams.GameServer{}, errors.New("no matching DS available to claim"))

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	// Create request
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-456",
		Namespace:       "mygame",
		Deployment:      "production-fleet",
		RequestedRegion: []string{"us-west-2"},
	}

	// Execute
	response, err := server.CreateGameSession(context.Background(), req)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to claim DS from AMS")
	assert.Contains(t, err.Error(), "no matching DS available to claim")

	mockAMSClient.AssertExpectations(t)
}

func TestSessionDSM_CreateGameSession_ValidationError(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	// Setup expected call to return validation error
	mockAMSClient.On("ClaimDS", mock.Anything, mock.Anything).
		Return(ams.GameServer{}, errors.New("validation error: namespace is required"))

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	// Create request with missing namespace
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-789",
		Namespace:       "",
		Deployment:      "production-fleet",
		RequestedRegion: []string{"us-west-2"},
	}

	// Execute
	response, err := server.CreateGameSession(context.Background(), req)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "namespace is required")

	mockAMSClient.AssertExpectations(t)
}

func TestSessionDSM_TerminateGameSession_Success(t *testing.T) {
	// Create server (no AMS client needed for termination)
	server := &SessionDSM{}

	// Create request
	req := &sessiondsm.RequestTerminateGameSession{
		SessionId: "session-123",
		Namespace: "mygame",
		Zone:      "us-west-2",
	}

	// Execute
	response, err := server.TerminateGameSession(context.Background(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "session-123", response.SessionId)
	assert.Equal(t, "mygame", response.Namespace)
	assert.True(t, response.Success)
	assert.Equal(t, "Session terminated successfully", response.Reason)
}

func TestSessionDSM_CreateGameSessionAsync_Success(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	// Note: We won't verify the actual async call in this test
	// as it runs in a goroutine. We just verify the immediate response.

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	// Create request
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-async-123",
		Namespace:       "mygame",
		Deployment:      "production-fleet",
		SessionData:     `{"mapName":"Desert"}`,
		RequestedRegion: []string{"us-west-2"},
		MaximumPlayer:   10,
		ClientVersion:   "1.0.0",
		GameMode:        "deathmatch",
	}

	// Execute
	response, err := server.CreateGameSessionAsync(context.Background(), req)

	// Verify immediate response
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Success)
	assert.Equal(t, "Game session creation initiated", response.Message)

	// Note: The goroutine will execute asynchronously, so we can't verify
	// the ClaimDS call here without complex synchronization
}

func TestSessionDSM_CreateGameSession_MultipleRegions(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	// Setup expected call with multiple regions
	expectedGameServer := ams.GameServer{
		SessionID:  "session-multi",
		Namespace:  "mygame",
		PodName:    "server-xyz",
		IP:         "10.0.0.1",
		Port:       8080,
		Region:     "us-east-1", // AMS chose east over west
		Deployment: "multi-region-fleet",
		StatusText: "claimed",
	}

	mockAMSClient.On("ClaimDS", mock.Anything, mock.MatchedBy(func(params ams.DsClaimParameters) bool {
		return len(params.Regions) == 2 &&
			params.Regions[0] == "us-west-2" &&
			params.Regions[1] == "us-east-1"
	})).Return(expectedGameServer, nil)

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	// Create request with multiple regions
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-multi",
		Namespace:       "mygame",
		Deployment:      "multi-region-fleet",
		RequestedRegion: []string{"us-west-2", "us-east-1"},
		MaximumPlayer:   20,
	}

	// Execute
	response, err := server.CreateGameSession(context.Background(), req)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "us-east-1", response.Region)
	assert.Equal(t, "us-east-1", response.CreatedRegion)

	mockAMSClient.AssertExpectations(t)
}

func TestSessionDSM_CreateGameSession_PreservesSessionData(t *testing.T) {
	// Create mock AMS client
	mockAMSClient := new(MockAMSClient)

	expectedGameServer := ams.GameServer{
		SessionID:  "session-data",
		Namespace:  "mygame",
		PodName:    "server-123",
		IP:         "192.168.1.1",
		Port:       7777,
		Region:     "us-west-2",
		Deployment: "test-fleet",
		StatusText: "claimed",
	}

	mockAMSClient.On("ClaimDS", mock.Anything, mock.Anything).
		Return(expectedGameServer, nil)

	// Create server
	server := &SessionDSM{
		AMSClient: mockAMSClient,
	}

	sessionData := `{"mapName":"CityMap","mode":"conquest","maxScore":1000}`

	// Create request
	req := &sessiondsm.RequestCreateGameSession{
		SessionId:       "session-data",
		Namespace:       "mygame",
		Deployment:      "test-fleet",
		SessionData:     sessionData,
		RequestedRegion: []string{"us-west-2"},
	}

	// Execute
	response, err := server.CreateGameSession(context.Background(), req)

	// Verify session data is preserved
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, sessionData, response.SessionData)

	mockAMSClient.AssertExpectations(t)
}
