package middlewares

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TokenMiddleware() grpc.UnaryServerInterceptor {
	token := os.Getenv("DRIVE_SECRET_TOKEN")

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if len(md["token"]) == 0 || md["token"][0] != token {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token")
			}
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "token not found")
		}

		return handler(ctx, req)
	}
}
