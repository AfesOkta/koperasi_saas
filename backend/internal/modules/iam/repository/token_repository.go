package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Register(ctx context.Context, token *model.UserDeviceToken) error
	GetByUserID(ctx context.Context, userID uint) ([]model.UserDeviceToken, error)
	DeleteByToken(ctx context.Context, deviceToken string) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Register(ctx context.Context, token *model.UserDeviceToken) error {
	// Use Save to handle upsert (unique index on device_token)
	return r.db.WithContext(ctx).Save(token).Error
}

func (r *tokenRepository) GetByUserID(ctx context.Context, userID uint) ([]model.UserDeviceToken, error) {
	var tokens []model.UserDeviceToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (r *tokenRepository) DeleteByToken(ctx context.Context, deviceToken string) error {
	return r.db.WithContext(ctx).Where("device_token = ?", deviceToken).Delete(&model.UserDeviceToken{}).Error
}
