package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	iamModel "github.com/koperasi-gresik/backend/internal/modules/iam/model"
	iamRepo "github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	"github.com/koperasi-gresik/backend/internal/modules/organization/dto"
	"github.com/koperasi-gresik/backend/internal/modules/organization/model"
	"github.com/koperasi-gresik/backend/internal/modules/organization/repository"
	"github.com/koperasi-gresik/backend/internal/shared/database/seeds"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrganizationService interface {
	Create(ctx context.Context, req dto.OrganizationCreateRequest) (*dto.OrganizationResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.OrganizationResponse, error)
	List(ctx context.Context, params pagination.Params) ([]dto.OrganizationResponse, int64, error)
	UpdateSettings(ctx context.Context, id uint, settings map[string]interface{}) (*dto.OrganizationResponse, error)
	Onboard(ctx context.Context, req dto.OnboardingRequest) (*dto.OnboardingResponse, error)
}

type organizationService struct {
	repo     repository.OrganizationRepository
	userRepo iamRepo.UserRepository
	roleRepo iamRepo.RoleRepository
	db       *gorm.DB
	rdb      *redis.Client
}

func NewOrganizationService(repo repository.OrganizationRepository, userRepo iamRepo.UserRepository, roleRepo iamRepo.RoleRepository, db *gorm.DB, rdb *redis.Client) OrganizationService {
	return &organizationService{
		repo:     repo,
		userRepo: userRepo,
		roleRepo: roleRepo,
		db:       db,
		rdb:      rdb,
	}
}

func (s *organizationService) Create(ctx context.Context, req dto.OrganizationCreateRequest) (*dto.OrganizationResponse, error) {
	slug := utils.GenerateUUID()[:8] // Basic slug for MVP, can be improved

	org := &model.Organization{
		Name:    req.Name,
		Slug:    slug,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err := s.repo.Create(ctx, org); err != nil {
		return nil, errors.New("failed to create organization")
	}

	return s.mapToResponse(org), nil
}

func (s *organizationService) GetByID(ctx context.Context, id uint) (*dto.OrganizationResponse, error) {
	// 1. Try Redis cache (24 Hour TTL)
	cacheKey := fmt.Sprintf("org:cfg:%d", id)
	if val, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil && val != "" {
		var cachedRes dto.OrganizationResponse
		if err := json.Unmarshal([]byte(val), &cachedRes); err == nil {
			return &cachedRes, nil
		}
	}

	// 2. Fetch from DB
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	res := s.mapToResponse(org)

	// 3. Cache the result
	if jsonStr, err := json.Marshal(res); err == nil {
		s.rdb.Set(ctx, cacheKey, jsonStr, 24*time.Hour)
	}

	return res, nil
}

func (s *organizationService) List(ctx context.Context, params pagination.Params) ([]dto.OrganizationResponse, int64, error) {
	orgs, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, 0, errors.New("failed to list organizations")
	}

	responses := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = *s.mapToResponse(&org)
	}

	return responses, total, nil
}

func (s *organizationService) UpdateSettings(ctx context.Context, id uint, settings map[string]interface{}) (*dto.OrganizationResponse, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("organization not found")
	}

	// In GORM with datatypes.JSON, we might need to handle the conversion
	// For simplicity, let's assume we can marshal it
	settingsBytes, _ := json.Marshal(settings)
	org.Settings = settingsBytes

	if err := s.repo.Update(ctx, org); err != nil {
		return nil, errors.New("failed to update organization settings")
	}

	res := s.mapToResponse(org)

	// Invalidate Cache
	cacheKey := fmt.Sprintf("org:cfg:%d", id)
	s.rdb.Del(ctx, cacheKey)

	return res, nil
}

func (s *organizationService) mapToResponse(org *model.Organization) *dto.OrganizationResponse {
	var settings map[string]interface{}
	if len(org.Settings) > 0 {
		json.Unmarshal(org.Settings, &settings)
	}

	return &dto.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Slug:      org.Slug,
		Email:     org.Email,
		Phone:     org.Phone,
		Address:   org.Address,
		Logo:      org.Logo,
		Plan:      org.Plan,
		Settings:  settings,
		Status:    org.Status,
		CreatedAt: org.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *organizationService) Onboard(ctx context.Context, req dto.OnboardingRequest) (*dto.OnboardingResponse, error) {
	var org *model.Organization
	var user *iamModel.User

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Organization
		slug := utils.GenerateUUID()[:8]
		org = &model.Organization{
			Name:    req.OrganizationName,
			Slug:    slug,
			Email:   req.Email,
			Phone:   req.Phone,
			Address: req.Address,
		}

		if err := tx.Create(org).Error; err != nil {
			return err
		}

		// 2. Hash Password
		hash, err := utils.HashPassword(req.AdminPassword)
		if err != nil {
			return err
		}

		// 3. Create Admin User
		user = &iamModel.User{
			Name:         req.AdminName,
			Email:        req.AdminEmail,
			PasswordHash: hash,
		}
		user.OrganizationID = org.ID

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// 4. Seed system roles for the new organization
		if err := seeds.SeedSystemRoles(ctx, tx, org.ID); err != nil {
			return err
		}

		// 5. Seed SAK ETAP Chart of Accounts
		if err := seeds.SeedCOASAKETAP(ctx, tx, org.ID); err != nil {
			return err
		}

		// 5. Assign admin role to the admin user
		var adminRole iamModel.Role
		if err := tx.Where("organization_id = ? AND name = ?", org.ID, "admin").First(&adminRole).Error; err == nil {
			tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT DO NOTHING", user.ID, adminRole.ID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.OnboardingResponse{
		Organization: *s.mapToResponse(org),
		AdminUser: fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}, nil
}
