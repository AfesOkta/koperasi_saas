package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/shu/model"
	"gorm.io/gorm"
)

type SHURepository interface {
	SaveConfig(ctx context.Context, config *model.SHUConfig) error
	ListConfigs(ctx context.Context, orgID uint) ([]model.SHUConfig, error)
	GetConfig(ctx context.Context, id uint) (*model.SHUConfig, error)
	CreateDistributions(ctx context.Context, dists []model.SHUDistribution) error
	GetMemberDistributions(ctx context.Context, orgID, memberID uint) ([]model.SHUDistribution, error)
}

type shuRepository struct {
	db *gorm.DB
}

func NewSHURepository(db *gorm.DB) SHURepository {
	return &shuRepository{db: db}
}

func (r *shuRepository) SaveConfig(ctx context.Context, config *model.SHUConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *shuRepository) ListConfigs(ctx context.Context, orgID uint) ([]model.SHUConfig, error) {
	var configs []model.SHUConfig
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Find(&configs).Error
	return configs, err
}

func (r *shuRepository) GetConfig(ctx context.Context, id uint) (*model.SHUConfig, error) {
	var config model.SHUConfig
	err := r.db.WithContext(ctx).First(&config, id).Error
	return &config, err
}

func (r *shuRepository) CreateDistributions(ctx context.Context, dists []model.SHUDistribution) error {
	return r.db.WithContext(ctx).Create(&dists).Error
}

func (r *shuRepository) GetMemberDistributions(ctx context.Context, orgID, memberID uint) ([]model.SHUDistribution, error) {
	var dists []model.SHUDistribution
	err := r.db.WithContext(ctx).Where("organization_id = ? AND member_id = ?", orgID, memberID).Find(&dists).Error
	return dists, err
}
