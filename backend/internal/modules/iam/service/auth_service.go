package service

import (
	"context"
	"errors"

	"github.com/koperasi-gresik/backend/internal/modules/iam/dto"
	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	jwtSecret string
	expHours  int
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, secret string, expHours int) AuthService {
	return &authService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		jwtSecret: secret,
		expHours:  expHours,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Look up user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != "active" {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	// Determine effective role (assuming first role for MVP)
	var roleID uint
	if len(user.Roles) > 0 {
		roleID = user.Roles[0].ID
	}

	// Generate Token
	token, err := utils.GenerateToken(user.ID, user.OrganizationID, roleID, user.Email, s.jwtSecret, s.expHours)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login (ignore errors natively here or perform async)

	return &dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: "refresh-token-placeholder", // Add refresh token logic later
		ExpiresIn:    s.expHours * 3600,
		User: dto.UserResponse{
			ID:             user.ID,
			OrganizationID: user.OrganizationID,
			Name:           user.Name,
			Email:          user.Email,
			Status:         user.Status,
		},
	}, nil
}

func (s *authService) Register(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Dummy logic for MVP: Assign default organization (ID 1)
	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Phone:        req.Phone,
	}
	user.OrganizationID = 1

	if req.RoleID > 0 {
		role, err := s.roleRepo.GetByID(ctx, user.OrganizationID, req.RoleID)
		if err == nil {
			user.Roles = []model.Role{*role}
		}
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user account")
	}

	return &dto.UserResponse{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Email:          user.Email,
		Status:         user.Status,
	}, nil
}
