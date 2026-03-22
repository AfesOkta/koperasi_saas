package service

import (
	"context"
	"errors"

	iamModel "github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"github.com/koperasi-gresik/backend/internal/modules/iam/repository"
	"github.com/koperasi-gresik/backend/internal/shared/middleware"
	"gorm.io/gorm"
)

// RoleService handles RBAC role management for an organization.
type RoleService interface {
	ListRoles(ctx context.Context, orgID uint) ([]iamModel.Role, error)
	GetRole(ctx context.Context, orgID, roleID uint) (*iamModel.Role, error)
	CreateRole(ctx context.Context, orgID uint, name, description string) (*iamModel.Role, error)
	UpdateRole(ctx context.Context, orgID, roleID uint, name, description string) (*iamModel.Role, error)
	DeleteRole(ctx context.Context, orgID, roleID uint) error
	AssignPermissions(ctx context.Context, orgID, roleID uint, permissionNames []string) error
	RemovePermissions(ctx context.Context, orgID, roleID uint, permissionNames []string) error
	AssignRoleToUser(ctx context.Context, orgID, userID, roleID uint) error
	ListPermissions(ctx context.Context) ([]iamModel.Permission, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	db       *gorm.DB
	cache    *middleware.PermissionCache
}

// NewRoleService creates a new RoleService.
func NewRoleService(roleRepo repository.RoleRepository, db *gorm.DB, cache *middleware.PermissionCache) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		db:       db,
		cache:    cache,
	}
}

func (s *roleService) ListRoles(ctx context.Context, orgID uint) ([]iamModel.Role, error) {
	return s.roleRepo.List(ctx, orgID)
}

func (s *roleService) GetRole(ctx context.Context, orgID, roleID uint) (*iamModel.Role, error) {
	return s.roleRepo.GetByID(ctx, orgID, roleID)
}

func (s *roleService) CreateRole(ctx context.Context, orgID uint, name, description string) (*iamModel.Role, error) {
	role := &iamModel.Role{
		Name:        name,
		Description: description,
		IsSystem:    false,
		Version:     1,
	}
	role.OrganizationID = orgID

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, errors.New("failed to create role")
	}
	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, orgID, roleID uint, name, description string) (*iamModel.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, errors.New("role not found")
	}
	if role.IsSystem {
		return nil, errors.New("system roles cannot be renamed")
	}

	role.Name = name
	role.Description = description
	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, errors.New("failed to update role")
	}
	return role, nil
}

func (s *roleService) DeleteRole(ctx context.Context, orgID, roleID uint) error {
	role, err := s.roleRepo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return errors.New("role not found")
	}
	if role.IsSystem {
		return errors.New("system roles cannot be deleted")
	}
	return s.db.WithContext(ctx).Delete(role).Error
}

func (s *roleService) AssignPermissions(ctx context.Context, orgID, roleID uint, permissionNames []string) error {
	role, err := s.roleRepo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return errors.New("role not found")
	}

	var perms []iamModel.Permission
	if err := s.db.WithContext(ctx).Where("name IN ?", permissionNames).Find(&perms).Error; err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Append(perms); err != nil {
		return err
	}

	// Bump version → self-invalidates old Redis keys
	return s.bumpVersion(ctx, role)
}

func (s *roleService) RemovePermissions(ctx context.Context, orgID, roleID uint, permissionNames []string) error {
	role, err := s.roleRepo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return errors.New("role not found")
	}

	var perms []iamModel.Permission
	if err := s.db.WithContext(ctx).Where("name IN ?", permissionNames).Find(&perms).Error; err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Model(role).Association("Permissions").Delete(perms); err != nil {
		return err
	}

	return s.bumpVersion(ctx, role)
}

func (s *roleService) AssignRoleToUser(ctx context.Context, orgID, userID, roleID uint) error {
	// Validate role belongs to org
	if _, err := s.roleRepo.GetByID(ctx, orgID, roleID); err != nil {
		return errors.New("role not found in organization")
	}
	// Replace existing role (single role per user)
	return s.db.WithContext(ctx).Exec(
		"DELETE FROM user_roles WHERE user_id = ?; INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)",
		userID, userID, roleID,
	).Error
}

func (s *roleService) ListPermissions(ctx context.Context) ([]iamModel.Permission, error) {
	var perms []iamModel.Permission
	if err := s.db.WithContext(ctx).Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

// bumpVersion increments role.version, causing the old Redis cache key to be abandoned.
func (s *roleService) bumpVersion(ctx context.Context, role *iamModel.Role) error {
	return s.db.WithContext(ctx).
		Model(role).
		UpdateColumn("version", gorm.Expr("version + 1")).
		Error
}
