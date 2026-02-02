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
	FindAllIncludingInactive(ctx context.Context) ([]*Tenant, error)
	FindByIDIncludingInactive(ctx context.Context, id string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
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

// FindByID finds active tenant by ID (used by middleware)
func (r *gormRepository) FindByID(ctx context.Context, id string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of error for not found
		}
		return nil, err
	}
	return &tenant, nil
}

// FindByIDIncludingInactive finds tenant by ID including inactive ones (for admin)
func (r *gormRepository) FindByIDIncludingInactive(ctx context.Context, id string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tenant, nil
}

// FindBySlug finds active tenant by slug (used by middleware)
func (r *gormRepository) FindBySlug(ctx context.Context, slug string) (*Tenant, error) {
	var tenant Tenant
	err := r.db.WithContext(ctx).Where("slug = ? AND is_active = ?", slug, true).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil instead of error
		}
		return nil, err
	}
	return &tenant, nil
}

// FindAll returns all active tenants
func (r *gormRepository) FindAll(ctx context.Context) ([]*Tenant, error) {
	var tenants []*Tenant
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&tenants).Error
	return tenants, err
}

// FindAllIncludingInactive returns all tenants including inactive (for admin)
func (r *gormRepository) FindAllIncludingInactive(ctx context.Context) ([]*Tenant, error) {
	var tenants []*Tenant
	err := r.db.WithContext(ctx).Find(&tenants).Error
	return tenants, err
}

func (r *gormRepository) Update(ctx context.Context, tenant *Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

// Delete performs soft delete by setting is_active = false
func (r *gormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&Tenant{}).Where("id = ?", id).Update("is_active", false).Error
}

// HardDelete permanently removes tenant from database
func (r *gormRepository) HardDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Tenant{}, "id = ?", id).Error
}
