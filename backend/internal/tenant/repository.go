package tenant

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository defines the interface for tenant data access
type Repository interface {
	Create(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, id string) (*Tenant, error)
	FindBySlug(ctx context.Context, slug string) (*Tenant, error)
	FindAll(ctx context.Context) ([]*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, tenant *Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *gormRepository) FindByID(ctx context.Context, id string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *gormRepository) FindBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.WithContext(ctx).Where("slug = ? AND is_active = ?", slug, true).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return &tenant, nil
}

func (r *gormRepository) FindAll(ctx context.Context) ([]*Tenant, error) {
	var tenants []*Tenant
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&tenants).Error
	return tenants, err
}

func (r *gormRepository) Update(ctx context.Context, tenant *Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *gormRepository) Delete(ctx context.Context, id string) error {
	// Soft delete by setting is_active = false
	return r.db.WithContext(ctx).Model(&Tenant{}).Where("id = ?", id).Update("is_active", false).Error
}
