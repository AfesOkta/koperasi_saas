package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"github.com/koperasi-gresik/backend/internal/shared/pagination"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, orgID, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByEmailAndOrg(ctx context.Context, orgID uint, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	List(ctx context.Context, orgID uint, params pagination.Params) ([]model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, orgID, id uint) (*model.User, error) {
	var user model.User
	query := r.db.WithContext(ctx).Preload("Roles")
	if orgID > 0 {
		query = query.Where("organization_id = ?", orgID)
	}
	if err := query.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailAndOrg(ctx context.Context, orgID uint, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Roles").Where("organization_id = ? AND email = ?", orgID, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	// Full save including associations if loaded
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) List(ctx context.Context, orgID uint, params pagination.Params) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Where("organization_id = ?", orgID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Roles").Scopes(params.Scope()).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
