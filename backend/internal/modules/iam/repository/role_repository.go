package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, orgID, id uint) (*model.Role, error)
	List(ctx context.Context, orgID uint) ([]model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	HasPermission(ctx context.Context, roleID uint, permissionName string) (bool, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, orgID, id uint) (*model.Role, error) {
	var role model.Role
	query := r.db.WithContext(ctx).Preload("Permissions")
	if orgID > 0 {
		query = query.Where("organization_id = ?", orgID)
	}
	if err := query.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context, orgID uint) ([]model.Role, error) {
	var roles []model.Role
	query := r.db.WithContext(ctx).Where("organization_id = ?", orgID)

	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) HasPermission(ctx context.Context, roleID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("role_permissions").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.name = ?", roleID, permissionName).
		Count(&count).Error
	return count > 0, err
}
