package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/billing/model"
	"gorm.io/gorm"
)

type BillingRepository interface {
	GetSubscription(ctx context.Context, orgID uint) (*model.OrgSubscription, error)
	ListPlans(ctx context.Context) ([]model.SubscriptionPlan, error)
	GetPlanByID(ctx context.Context, id uint) (*model.SubscriptionPlan, error)
	CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error
	UpdatePlan(ctx context.Context, plan *model.SubscriptionPlan) error
	DeletePlan(ctx context.Context, id uint) error
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

func (r *billingRepository) ListPlans(ctx context.Context) ([]model.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan
	err := r.db.WithContext(ctx).Find(&plans).Error
	return plans, err
}

func (r *billingRepository) GetPlanByID(ctx context.Context, id uint) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	err := r.db.WithContext(ctx).First(&plan, id).Error
	return &plan, err
}

func (r *billingRepository) CreatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *billingRepository) UpdatePlan(ctx context.Context, plan *model.SubscriptionPlan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *billingRepository) DeletePlan(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.SubscriptionPlan{}, id).Error
}

func (r *billingRepository) GetPlanByCode(ctx context.Context, code string) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&plan).Error
	return &plan, err
}

func (r *billingRepository) UpdateSubscription(ctx context.Context, sub *model.OrgSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}
