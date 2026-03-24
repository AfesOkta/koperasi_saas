package service

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/billing/dto"
	"github.com/koperasi-gresik/backend/internal/modules/billing/model"
	"github.com/koperasi-gresik/backend/internal/modules/billing/repository"
)

type BillingService interface {
	GetSubscription(ctx context.Context, orgID uint) (*dto.OrgSubscriptionResponse, error)
	ListPlans(ctx context.Context) ([]dto.SubscriptionPlanResponse, error)
	GetPlanByID(ctx context.Context, id uint) (*dto.SubscriptionPlanResponse, error)
	CreatePlan(ctx context.Context, req dto.SubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error)
	UpdatePlan(ctx context.Context, id uint, req dto.SubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error)
	DeletePlan(ctx context.Context, id uint) error
}

type billingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingService{repo: repo}
}

func (s *billingService) GetSubscription(ctx context.Context, orgID uint) (*dto.OrgSubscriptionResponse, error) {
	sub, err := s.repo.GetSubscription(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return &dto.OrgSubscriptionResponse{
		ID:             sub.ID,
		OrganizationID: sub.OrganizationID,
		PlanID:         sub.PlanID,
		Plan:           *s.mapPlanToResponse(&sub.Plan),
		StartDate:      sub.StartDate.Format("2006-01-02"),
		EndDate:        sub.EndDate.Format("2006-01-02"),
		Status:         sub.Status,
	}, nil
}

func (s *billingService) ListPlans(ctx context.Context) ([]dto.SubscriptionPlanResponse, error) {
	plans, err := s.repo.ListPlans(ctx)
	if err != nil {
		return nil, err
	}

	var res []dto.SubscriptionPlanResponse
	for _, p := range plans {
		res = append(res, *s.mapPlanToResponse(&p))
	}
	return res, nil
}

func (s *billingService) GetPlanByID(ctx context.Context, id uint) (*dto.SubscriptionPlanResponse, error) {
	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapPlanToResponse(plan), nil
}

func (s *billingService) CreatePlan(ctx context.Context, req dto.SubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error) {
	plan := &model.SubscriptionPlan{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Price:       req.Price,
		MaxUsers:    req.MaxUsers,
		MaxMembers:  req.MaxMembers,
	}

	if err := s.repo.CreatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return s.mapPlanToResponse(plan), nil
}

func (s *billingService) UpdatePlan(ctx context.Context, id uint, req dto.SubscriptionPlanRequest) (*dto.SubscriptionPlanResponse, error) {
	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	plan.Name = req.Name
	plan.Code = req.Code
	plan.Description = req.Description
	plan.Price = req.Price
	plan.MaxUsers = req.MaxUsers
	plan.MaxMembers = req.MaxMembers

	if err := s.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return s.mapPlanToResponse(plan), nil
}

func (s *billingService) DeletePlan(ctx context.Context, id uint) error {
	return s.repo.DeletePlan(ctx, id)
}

func (s *billingService) mapPlanToResponse(p *model.SubscriptionPlan) *dto.SubscriptionPlanResponse {
	return &dto.SubscriptionPlanResponse{
		ID:          p.ID,
		Name:        p.Name,
		Code:        p.Code,
		Description: p.Description,
		Price:       p.Price,
		MaxUsers:    p.MaxUsers,
		MaxMembers:  p.MaxMembers,
	}
}
