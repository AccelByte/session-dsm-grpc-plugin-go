# AMS Session DSM Implementation

This package implements the Session DSM gRPC plugin for AccelByte Multiplayer Servers (AMS).

## Overview

The AMS implementation allows the Session Service to claim dedicated servers from AMS fleets dynamically based on session requirements.

## Features

- **Synchronous DS Claim** (`CreateGameSession`): Claims a DS immediately and returns the server details
- **Asynchronous DS Claim** (`CreateGameSessionAsync`): Claims a DS in the background and updates the session via callback
- **Session Termination** (`TerminateGameSession`): Handles session cleanup (AMS automatically recycles servers)

## How It Works

### 1. CreateGameSession (Synchronous)

```
Session Service → gRPC CreateGameSession → AMS Plugin
                                          ↓
                                    Claim DS from AMS Fleet
                                          ↓
                                    Return Server Details
```

**Flow:**
1. Session service calls `CreateGameSession` with session details
2. Plugin validates claim parameters (namespace, regions, claim keys)
3. Plugin calls AMS API to claim a server from the fleet
4. AMS returns server details (IP, port, serverID, region)
5. Plugin returns server info to session service
6. Session becomes AVAILABLE with DS connection details

### 2. CreateGameSessionAsync (Asynchronous)

```
Session Service → gRPC CreateGameSessionAsync → AMS Plugin
                ← Response: "Initiated"

                                          Background Worker
                                          ↓
                                    Claim DS from AMS Fleet
                                          ↓
                                    Update Session Service
                                    via AdminUpdateDSInformation API
```

**Flow:**
1. Session service calls `CreateGameSessionAsync`
2. Plugin immediately returns success response
3. Background goroutine claims DS from AMS
4. On success: Updates session with `AdminUpdateDSInformation` API
5. On failure: Updates session with error status
6. Session service receives webhook notification about DS status

### 3. TerminateGameSession

**Flow:**
1. Session service calls `TerminateGameSession`
2. Plugin returns success (no actual AMS API call needed)
3. AMS automatically recycles the server back to the fleet pool

## Configuration

### Required Parameters

**From `RequestCreateGameSession`:**
- `session_id`: Unique session identifier
- `namespace`: Game namespace
- `deployment`: Used as claim key for fleet matching
- `requested_region`: List of preferred regions (prioritized by AMS)
- `client_version`: Game client version (optional, for filtering)
- `game_mode`: Game mode (optional, for filtering)

### AMS Client Setup

```go
import (
    "session-dsm-grpc-plugin/pkg/client/ams"
    "session-dsm-grpc-plugin/pkg/server/ams"
)

// Create AMS client
amsClient := ams.New(
    amsBaseURL,              // e.g., "https://ams.example.com"
    func() string {          // Token provider function
        return iamClient.ClientToken()
    },
    30 * time.Second,        // Request timeout
)

// Create AMS Session DSM server
amsServer := &server.SessionDSM{
    AMSClient:     amsClient,
    SessionClient: sessionClient,
    IamClient:     iamClient,
}

// Register with gRPC
sessiondsm.RegisterSessionDsmServer(grpcServer, amsServer)
```

## AMS API Endpoints Used

### 1. Claim DS
```
PUT /v1/namespaces/{namespace}/servers/claim
{
  "claimKeys": ["deployment-name"],
  "regions": ["us-west-2", "us-east-1"],
  "sessionId": "session-123"
}

Response:
{
  "serverId": "server-abc123",
  "ip": "192.168.1.100",
  "ports": {
    "default": 7777,
    "query": 7778
  },
  "region": "us-west-2",
  "fleetId": "fleet-xyz"
}
```

## Error Handling

### Claim Failures

**No Available Servers (404):**
```json
{
  "errorMessage": "no matching DS available to claim",
  "traceId": "trace-123"
}
```
- Plugin returns error to session service
- Session status: `FAILED_TO_REQUEST`

**Invalid Parameters (400):**
```json
{
  "errorMessage": "malformed request: missing claimKeys",
  "traceId": "trace-456"
}
```

### Retry Logic

- Exponential backoff with max 5 retries
- Client errors (4xx): No retry (permanent failure)
- Server errors (5xx): No retry (permanent failure)
- Network errors: Retried up to 5 times

## Comparison with Other Implementations

| Feature | AMS | DSMC | Custom |
|---------|-----|------|--------|
| **Provisioning** | Claim from pool | Create on-demand | Manual |
| **Speed** | Fast (pre-warmed) | Slower (cold start) | Varies |
| **Scaling** | Auto-scaling fleet | On-demand | Manual |
| **Cost** | Pool-based | Per-instance | Manual |
| **API** | Claim/Return | Create/Terminate | Custom |

## Integration with Session Service

### Session Configuration

```json
{
  "type": "DS",
  "dsSource": "AMS",
  "deployment": "my-game-fleet",
  "requestedRegions": ["us-west-2", "us-east-1"],
  "preferredClaimKeys": ["production-fleet"],
  "fallbackClaimKeys": ["staging-fleet"],
  "amsClaimTimeoutMinutes": 2
}
```

### Session Lifecycle

```
1. Session Created (status: NEED_TO_REQUEST)
2. gRPC Plugin Called (CreateGameSession/Async)
3. DS Claimed from AMS
4. Session Updated (status: AVAILABLE)
5. Players Connect
6. Session Ends
7. TerminateGameSession Called
8. Server Returned to Pool Automatically
```

## Monitoring & Logs

### Log Fields

All operations log the following fields:
- `sessionID`: Session identifier
- `namespace`: Game namespace
- `regions`: Requested regions
- `claimKeys`: Fleet claim keys
- `serverID`: Claimed server ID
- `ip`: Server IP address
- `port`: Server port
- `region`: Server region
- `fleetID`: AMS fleet ID

### Example Logs

**Successful Claim:**
```json
{
  "level": "info",
  "msg": "Successfully claimed DS from AMS",
  "sessionID": "abc123",
  "serverID": "server-xyz",
  "ip": "192.168.1.100",
  "port": 7777,
  "region": "us-west-2"
}
```

**Failed Claim:**
```json
{
  "level": "error",
  "msg": "no matching DS available to claim",
  "sessionID": "abc123",
  "namespace": "mygame",
  "statusCode": 404
}
```

## Testing

See `client_test.go` for unit tests covering:
- Parameter validation
- Successful claims
- Error handling
- Retry logic
- Response parsing

## References

- [Session Service AMS Integration](/Users/desnia/go/justice-session-service/pkg/ams/client.go)
- [Session DSM Proto Definition](/Users/desnia/go/github/session-dsm-grpc-plugin-go/pkg/proto/session-dsm.proto)
- [AMS Documentation](https://docs.accelbyte.io/gaming-services/services/play/ams/)
