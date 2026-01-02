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

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	token, user, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		Token:    token,
		UserId:   user.ID,
		TenantId: user.TenantID,
		Role:     user.Role,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authv1.TokenRequest) (*authv1.TokenResponse, error) {
	valid, user := h.service.ValidateToken(req.Token)
	if !valid {
		return &authv1.TokenResponse{Valid: false}, nil
	}

	return &authv1.TokenResponse{
		Valid:    true,
		UserId:   user.ID,
		TenantId: user.TenantID,
		Role:     user.Role,
	}, nil
}
