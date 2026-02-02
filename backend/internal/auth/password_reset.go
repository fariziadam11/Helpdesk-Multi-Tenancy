package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/user"
)

// RequestPasswordReset handles password reset request
func (s *service) RequestPasswordReset(ctx context.Context, tenantID, email string) error {
	// Find user by email within tenant
	u, err := s.userRepo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user by email")
		return err
	}

	// If user not found, silently return success (don't reveal if email exists)
	if u == nil {
		s.logger.Infof("password reset requested for non-existent email: %s", email)
		return nil
	}

	// Generate reset token
	token, err := GenerateResetToken()
	if err != nil {
		s.logger.WithError(err).Error("failed to generate reset token")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to generate reset token", err)
	}

	// Create reset token record with 1 hour expiration
	resetToken := &user.ResetToken{
		TenantID:  tenantID,
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.userRepo.CreateResetToken(ctx, tenantID, resetToken); err != nil {
		s.logger.WithError(err).Error("failed to create reset token")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to create reset token", err)
	}

	// Build reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, token)

	// Send email
	if err := s.emailClient.SendPasswordResetEmail(u.Email, resetLink); err != nil {
		s.logger.WithError(err).Error("failed to send password reset email")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to send password reset email", err)
	}

	s.logger.Infof("password reset email sent to: %s", u.Email)
	return nil
}

// ResetPassword handles password reset with token
func (s *service) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Get reset token (token is globally unique, contains tenantID)
	resetToken, err := s.userRepo.GetResetToken(ctx, token)
	if err != nil {
		s.logger.WithError(err).Error("failed to get reset token")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to get reset token", err)
	}

	if resetToken == nil {
		return errors.NewAppError(errors.ErrCodeInvalidInput, "invalid or expired reset token", nil)
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		// Delete expired token
		_ = s.userRepo.DeleteResetToken(ctx, token)
		return errors.NewAppError(errors.ErrCodeInvalidInput, "reset token has expired", nil)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithError(err).Error("failed to hash password")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to hash password", err)
	}

	// Update user password within tenant
	if err := s.userRepo.UpdatePassword(ctx, resetToken.TenantID, resetToken.UserID, string(hashedPassword)); err != nil {
		s.logger.WithError(err).Error("failed to update password")
		return errors.NewAppError(errors.ErrCodeInternal, "failed to update password", err)
	}

	// Delete used token
	if err := s.userRepo.DeleteResetToken(ctx, token); err != nil {
		s.logger.WithError(err).Warn("failed to delete reset token after use")
		// Don't return error here, password was already updated
	}

	s.logger.Infof("password reset successful for user: %s", resetToken.UserID)
	return nil
}
