package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/organization/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *model.Organization) error
	GetByID(ctx context.Context, id uint) (*model.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*model.Organization, error)
	Update(ctx context.Context, org *model.Organization) error
	List(ctx context.Context, params pagination.Params) ([]model.Organization, int64, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uint) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) Update(ctx context.Context, org *model.Organization) error {
	return r.db.WithContext(ctx).Save(org).Error
}

func (r *organizationRepository) List(ctx context.Context, params pagination.Params) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Organization{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Scopes(params.Scope()).Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}
