package service

import (
	"context"
	"errors"
	"time" // Added for time.Now().Unix()

	"github.com/koperasi-gresik/backend/internal/modules/iam/dto"
	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	"github.com/koperasi-gresik/backend/internal/shared/event"
	"github.com/koperasi-gresik/backend/internal/shared/utils"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	Register(ctx context.Context, req dto.UserCreateRequest) (*dto.UserResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	RegisterDeviceToken(ctx context.Context, orgID, userID uint, req dto.DeviceTokenRequest) error
}

type authService struct {
	userRepo  repository.UserRepository
	roleRepo  repository.RoleRepository
	tokenRepo repository.TokenRepository
	publisher event.Publisher
	secret    string
	expHours  int
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenRepo repository.TokenRepository,
	publisher event.Publisher,
	secret string,
	expHours int,
) AuthService {
	return &authService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		tokenRepo: tokenRepo,
		publisher: publisher,
		secret:    secret,
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
	var roleVersion int
	if len(user.Roles) > 0 {
		roleID = user.Roles[0].ID
		roleVersion = user.Roles[0].Version
	}

	// Generate Access Token
	accessToken, err := utils.GenerateToken(user.ID, user.OrganizationID, roleID, roleVersion, user.Email, s.secret, s.expHours)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate Refresh Token (30 days for mobile/web persistent)
	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.secret, 720) 
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

	// Emit EventUserCreated
	go func() {
		evt := event.Event{
			Type:           event.EventUserCreated,
			AggregateID:    user.ID,
			OrganizationID: user.OrganizationID,
			Payload:        user,
			Timestamp:      time.Now().Unix(),
		}
		_ = s.publisher.Publish(context.Background(), evt)
	}()

	return &dto.UserResponse{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Email:          user.Email,
		Status:         user.Status,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	userID, err := utils.VerifyRefreshToken(refreshToken, s.secret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.GetByID(ctx, 0, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != "active" {
		return nil, errors.New("account is disabled")
	}

	// Determine effective role
	var roleID uint
	var roleVersion int
	if len(user.Roles) > 0 {
		roleID = user.Roles[0].ID
		roleVersion = user.Roles[0].Version
	}

	accessToken, err := utils.GenerateToken(user.ID, user.OrganizationID, roleID, roleVersion, user.Email, s.secret, s.expHours)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken, 
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

func (s *authService) RegisterDeviceToken(ctx context.Context, orgID, userID uint, req dto.DeviceTokenRequest) error {
	token := &model.UserDeviceToken{
		UserID:      userID,
		DeviceToken: req.DeviceToken,
		DeviceType:  req.DeviceType,
	}
	token.OrganizationID = orgID
	
	return s.tokenRepo.Register(ctx, token)
}
