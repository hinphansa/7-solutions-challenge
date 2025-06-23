package grpc

import (
	"context"
	"strings"

	"github.com/hinphansa/7-solutions-challenge/internal/adapters/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ctxKey string

const (
	userIDKey ctxKey = "user_id"
)

// UnaryAuthInterceptor is a gRPC middleware that handles JWT authentication
func UnaryAuthInterceptor(jwtManager *auth.JWTMaker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Skip authentication for public endpoints
		if isPublicEndpoint(info.FullMethod) {
			return handler(ctx, req)
		}

		// Get token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		// Extract token from "Bearer <token>"
		authHeader := values[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT
		id, valid, err := jwtManager.Verify(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}
		if !valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, userIDKey, id)
		return handler(ctx, req)
	}
}

// isPublicEndpoint checks if the endpoint requires authentication
func isPublicEndpoint(fullMethod string) bool {
	publicEndpoints := map[string]bool{
		"/user.UserService/CreateUser":  true,
		"/user.UserService/GetUserById": true,
		"/user.UserService/ListUsers":   true,
		"/user.UserService/Login":       true,
	}
	return publicEndpoints[fullMethod]
}
