package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/invgate"
	"werk-ticketing/internal/tenant"
	"werk-ticketing/internal/user"
)

// Service exposes authentication related use cases.
// All methods now require tenantID for multi-tenant support.
type Service interface {
	Register(ctx context.Context, tenantID string, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, tenantID string, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, tenantID, refreshToken string) (*AuthResponse, error)
	RevokeToken(ctx context.Context, token string) error
	ParseToken(token string) (*jwt.RegisteredClaims, error)
	IsTokenBlacklisted(token string) bool
	// Password reset methods
	RequestPasswordReset(ctx context.Context, tenantID, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	userRepo      user.Repository
	tenantRepo    tenant.Repository
	invgateClient invgate.Service
	jwtSecret     []byte
	blacklist     *TokenBlacklist
	logger        *logrus.Logger
	emailClient   EmailClient
	frontendURL   string
}

// NewService instantiates auth service.
func NewService(
	userRepo user.Repository,
	tenantRepo tenant.Repository,
	invgateClient invgate.Service,
	jwtSecret string,
	logger *logrus.Logger,
	emailClient EmailClient,
	frontendURL string,
) Service {
	return &service{
		userRepo:      userRepo,
		tenantRepo:    tenantRepo,
		invgateClient: invgateClient,
		jwtSecret:     []byte(jwtSecret),
		blacklist:     NewTokenBlacklist(),
		logger:        logger,
		emailClient:   emailClient,
		frontendURL:   frontendURL,
	}
}

// EmailClient interface for sending emails
type EmailClient interface {
	SendPasswordResetEmail(to, resetLink string) error
}
