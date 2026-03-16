package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	GetSubscription(ctx context.Context, orgID uint) (*model.OrgSubscription, error)
	CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error
	GetPlanByCode(ctx context.Context, code string) (*model.SubscriptionPlan, error)
	UpdateSubscription(ctx context.Context, sub *model.OrgSubscription) error
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db: db}
}

func (r *billingRepository) GetSubscription(ctx context.Context, orgID uint) (*model.OrgSubscription, error) {
	var sub model.OrgSubscription
	err := r.db.WithContext(ctx).Preload("Plan").Where("organization_id = ?", orgID).First(&sub).Error
	return &sub, err
}

func (r *billingRepository) CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *billingRepository) GetPlanByCode(ctx context.Context, code string) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&plan).Error
	return &plan, err
}

func (r *billingRepository) UpdateSubscription(ctx context.Context, sub *model.OrgSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}
