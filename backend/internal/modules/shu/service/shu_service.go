package service

import (
	"context"
	"errors"

	"github.com/koperasi-gresik/backend/internal/modules/shu/dto"
	"github.com/koperasi-gresik/backend/internal/modules/shu/model"
	"github.com/koperasi-gresik/backend/internal/modules/shu/repository"
)

type SHUService interface {
	CreateConfig(ctx context.Context, orgID uint, req dto.SHUConfigRequest) (*dto.SHUConfigResponse, error)
	ListConfigs(ctx context.Context, orgID uint) ([]dto.SHUConfigResponse, error)
	CalculateSHU(ctx context.Context, orgID, configID uint) error
	GetMemberDistributions(ctx context.Context, orgID, memberID uint) ([]dto.SHUDistributionResponse, error)
}

type shuService struct {
	repo repository.SHURepository
}

func NewSHUService(repo repository.SHURepository) SHUService {
	return &shuService{repo: repo}
}

func (s *shuService) CreateConfig(ctx context.Context, orgID uint, req dto.SHUConfigRequest) (*dto.SHUConfigResponse, error) {
	config := &model.SHUConfig{
		Year:              req.Year,
		TotalSHU:          req.TotalSHU,
		MemberSavingsPct:  req.MemberSavingsPct,
		MemberBusinessPct: req.MemberBusinessPct,
		Status:            "draft",
	}
	config.OrganizationID = orgID

	if err := s.repo.SaveConfig(ctx, config); err != nil {
		return nil, err
	}

	return s.mapConfigToResponse(config), nil
}

func (s *shuService) ListConfigs(ctx context.Context, orgID uint) ([]dto.SHUConfigResponse, error) {
	configs, err := s.repo.ListConfigs(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var res []dto.SHUConfigResponse
	for _, c := range configs {
		res = append(res, *s.mapConfigToResponse(&c))
	}
	return res, nil
}

func (s *shuService) CalculateSHU(ctx context.Context, orgID, configID uint) error {
	config, err := s.repo.GetConfig(ctx, configID)
	if err != nil {
		return err
	}
	if config.OrganizationID != orgID {
		return errors.New("unauthorized access to config")
	}

	if config.Status != "draft" {
		return errors.New("only draft config can be calculated")
	}

	// MVP Logic: Mock member distributions 
	// In reality: 
	// 1. Get all members' avg savings balances
	// 2. Get all members' business transactions
	// 3. Compute proportions

	savingsPool := config.TotalSHU * (config.MemberSavingsPct / 100)
	businessPool := config.TotalSHU * (config.MemberBusinessPct / 100)

	// Since we don't have a direct repo to member balances in SHU module for MVP,
	// we will create a dummy distribution for member ID 1.
	var dists []model.SHUDistribution
	dists = append(dists, model.SHUDistribution{
		SHUConfigID:   config.ID,
		MemberID:      1,
		SavingsShare:  savingsPool * 0.1,  // Mocking 10% share
		BusinessShare: businessPool * 0.1, // Mocking 10% share
		TotalAmount:   (savingsPool * 0.1) + (businessPool * 0.1),
		Status:        "pending",
	})
	
	for i := range dists {
		dists[i].OrganizationID = orgID
	}

	if err := s.repo.CreateDistributions(ctx, dists); err != nil {
		return err
	}

	config.Status = "calculated"
	return s.repo.SaveConfig(ctx, config)
}

func (s *shuService) GetMemberDistributions(ctx context.Context, orgID, memberID uint) ([]dto.SHUDistributionResponse, error) {
	dists, err := s.repo.GetMemberDistributions(ctx, orgID, memberID)
	if err != nil {
		return nil, err
	}

	var res []dto.SHUDistributionResponse
	for _, d := range dists {
		res = append(res, dto.SHUDistributionResponse{
			ID:            d.ID,
			MemberID:      d.MemberID,
			SavingsShare:  d.SavingsShare,
			BusinessShare: d.BusinessShare,
			TotalAmount:   d.TotalAmount,
			Status:        d.Status,
			CreatedAt:     d.CreatedAt,
		})
	}
	return res, nil
}

func (s *shuService) mapConfigToResponse(c *model.SHUConfig) *dto.SHUConfigResponse {
	return &dto.SHUConfigResponse{
		ID:                c.ID,
		Year:              c.Year,
		TotalSHU:          c.TotalSHU,
		MemberSavingsPct:  c.MemberSavingsPct,
		MemberBusinessPct: c.MemberBusinessPct,
		Status:            c.Status,
		CreatedAt:         c.CreatedAt,
	}
}
