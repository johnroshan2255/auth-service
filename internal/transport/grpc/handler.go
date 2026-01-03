// Package grpc provides gRPC handlers for backend-to-backend communication.
// This is used when the auth service needs to communicate with other backend services (e.g., core service).
// For frontend communication, use the HTTP handlers in internal/transport/http.
package grpc

import (
	"context"

	authv1 "github.com/johnroshan2255/auth-service/proto/auth/v1"
	"github.com/johnroshan2255/auth-service/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
	authv1.UnimplementedAuthServiceServer
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authv1.TokenRequest) (*authv1.TokenResponse, error) {
	valid, user := h.service.ValidateToken(req.Token)
	if !valid {
		return &authv1.TokenResponse{Valid: false}, nil
	}

	return &authv1.TokenResponse{
		Valid:    true,
		UserId:   user.UUID,
		TenantId: user.TenantID,
		Role:     user.Role,
	}, nil
}
