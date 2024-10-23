// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package common

import (
	"context"
	"os"
	"strings"

	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var OAuth *iam.OAuth20Service

func UnaryAuthServerIntercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !skipCheckAuthorizationMetadata(info.FullMethod) {
		err := checkAuthorizationMetadata(ctx)

		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func StreamAuthServerIntercept(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if !skipCheckAuthorizationMetadata(info.FullMethod) {
		err := checkAuthorizationMetadata(ss.Context())

		if err != nil {
			return err
		}
	}

	return handler(srv, ss)
}

func skipCheckAuthorizationMetadata(fullMethod string) bool {
	if strings.HasPrefix(fullMethod, "/grpc.reflection.v1alpha.ServerReflection/") {
		return true
	}

	if strings.HasPrefix(fullMethod, "/grpc.health.v1.Health/") {
		return true
	}

	return false
}

func checkAuthorizationMetadata(ctx context.Context) error {
	if OAuth == nil {
		return status.Error(codes.Internal, "authorization token validator is not set")
	}

	meta, found := metadata.FromIncomingContext(ctx)

	if !found {
		return status.Error(codes.Unauthenticated, "metadata is missing")
	}

	if _, ok := meta["authorization"]; !ok {
		return status.Error(codes.Unauthenticated, "authorization metadata is missing")
	}

	if len(meta["authorization"]) == 0 {
		return status.Error(codes.Unauthenticated, "authorization metadata length is 0")
	}

	authorization := meta["authorization"][0]
	token := strings.TrimPrefix(authorization, "Bearer ")
	namespace := os.Getenv("AB_NAMESPACE")

	err := OAuth.Validate(token, nil, &namespace, nil)

	if err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	return nil
}
