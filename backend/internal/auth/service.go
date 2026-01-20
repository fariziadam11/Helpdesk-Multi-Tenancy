package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/invgate"
	"werk-ticketing/internal/user"
)

// Service exposes authentication related use cases.
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	RevokeToken(ctx context.Context, token string) error
	ParseToken(token string) (*jwt.RegisteredClaims, error)
	IsTokenBlacklisted(token string) bool
	// Password reset methods
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type service struct {
	userRepo      user.Repository
	invgateClient invgate.Service
	jwtSecret     []byte
	blacklist     *TokenBlacklist
	logger        *logrus.Logger
	companyID     int
	groupID       int
	locationID    int
	emailClient   EmailClient
	frontendURL   string
}

// NewService instantiates auth service.
func NewService(repo user.Repository, invgateClient invgate.Service, jwtSecret string, logger *logrus.Logger, companyID, groupID, locationID int, emailClient EmailClient, frontendURL string) Service {
	return &service{
		userRepo:      repo,
		invgateClient: invgateClient,
		jwtSecret:     []byte(jwtSecret),
		blacklist:     NewTokenBlacklist(),
		logger:        logger,
		companyID:     companyID,
		groupID:       groupID,
		locationID:    locationID,
		emailClient:   emailClient,
		frontendURL:   frontendURL,
	}
}

// EmailClient interface for sending emails
type EmailClient interface {
	SendPasswordResetEmail(to, resetLink string) error
}
