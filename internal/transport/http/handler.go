package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johnroshan2255/auth-service/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserUUID   string `json:"user_uuid"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}

type SignupRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=3,max=50"`
	Password    string `json:"password" binding:"required,min=6"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
}

type SignupResponse struct {
	Token    string `json:"token"`
	UserUUID   string `json:"user_uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
}

type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type ValidateTokenResponse struct {
	Valid    bool   `json:"valid"`
	UserUUID   string `json:"user_uuid"`
	TenantID string `json:"tenant_id,omitempty"`
	Role     string `json:"role,omitempty"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:    token,
		UserUUID:   user.UUID,
		TenantID: user.TenantID,
		Role:     user.Role,
	})
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid, user := h.service.ValidateToken(req.Token)
	if !valid {
		c.JSON(http.StatusOK, ValidateTokenResponse{Valid: false})
		return
	}

	c.JSON(http.StatusOK, ValidateTokenResponse{
		Valid:    true,
		UserUUID:   user.UUID,
		TenantID: user.TenantID,
		Role:     user.Role,
	})
}

// GetCurrentUser returns the current authenticated user info
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	tenantID, _ := c.Get("tenant_id")
	role, _ := c.Get("role")

		c.JSON(http.StatusOK, gin.H{
		"user_uuid":   userUUID,
		"tenant_id": tenantID,
		"role":      role,
	})
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.service.Signup(
		c.Request.Context(),
		req.Email,
		req.Username,
		req.Password,
		req.PhoneNumber,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SignupResponse{
		Token:    token,
		UserUUID:   user.UUID,
		Email:    user.Email,
		Username: user.Username,
		TenantID: user.TenantID,
		Role:     user.Role,
	})
}

// HealthCheck checks the health of the auth service
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Auth service is running"})
}