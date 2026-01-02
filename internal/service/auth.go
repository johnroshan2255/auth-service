package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/johnroshan2255/auth-service/internal/model"
)

type AuthService struct {
	repo UserRepository
	// cache can be added here if needed
}

var jwtKey = []byte("replace_with_strong_secret")

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Login authenticates the user and returns a JWT
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

// ValidateToken parses JWT and returns user info
func (s *AuthService) ValidateToken(tokenStr string) (bool, *model.User) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	return true, &model.User{
		ID:       claims["user_id"].(string),
		TenantID: claims["tenant_id"].(string),
		Role:     claims["role"].(string),
	}
}
