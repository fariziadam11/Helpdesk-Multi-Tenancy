package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Repository abstracts data persistence for users.
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	// Password reset methods
	CreateResetToken(ctx context.Context, token *ResetToken) error
	GetResetToken(ctx context.Context, token string) (*ResetToken, error)
	DeleteResetToken(ctx context.Context, token string) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository builds a Gorm-backed user repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		// Check if error is due to duplicate key (unique constraint violation)
		if isDuplicateKeyError(err) {
			return &DuplicateKeyError{
				Field: "email",
				Value: user.Email,
				Err:   err,
			}
		}
		return err
	}
	return nil
}

// DuplicateKeyError represents a duplicate key constraint violation
type DuplicateKeyError struct {
	Field string
	Value string
	Err   error
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate key violation on field '%s' with value '%s'", e.Field, e.Value)
}

func (e *DuplicateKeyError) Unwrap() error {
	return e.Err
}

// isDuplicateKeyError checks if the error is a duplicate key constraint violation
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// MySQL duplicate key error patterns
	if strings.Contains(errStr, "duplicate entry") ||
		strings.Contains(errStr, "1062") || // MySQL error code for duplicate entry
		strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key") {
		return true
	}

	return false
}

func (r *gormRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *gormRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}

// CreateResetToken creates a new password reset token
func (r *gormRepository) CreateResetToken(ctx context.Context, token *ResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetResetToken retrieves a reset token by its value
func (r *gormRepository) GetResetToken(ctx context.Context, token string) (*ResetToken, error) {
	var rt ResetToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

// DeleteResetToken deletes a reset token
func (r *gormRepository) DeleteResetToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&ResetToken{}, "token = ?", token).Error
}

// UpdatePassword updates a user's password
func (r *gormRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
