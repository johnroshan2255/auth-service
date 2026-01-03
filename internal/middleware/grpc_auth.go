package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var serviceKey string

func SetServiceKey(key string) {
	serviceKey = key
}

func BackendAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if serviceKey == "" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	keys := md.Get("service-key")
	if len(keys) == 0 || keys[0] != serviceKey {
		return nil, status.Errorf(codes.Unauthenticated, "invalid service key")
	}

	return handler(ctx, req)
}
