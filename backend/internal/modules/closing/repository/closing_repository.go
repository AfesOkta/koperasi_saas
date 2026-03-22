package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/closing/model"
	"gorm.io/gorm"
)

type ClosingRepository interface {
	CreateLog(ctx context.Context, log *model.ClosingLog) error
	UpdateLog(ctx context.Context, log *model.ClosingLog) error
	GetLogByDate(ctx context.Context, orgID uint, date string, typ string) (*model.ClosingLog, error)
	IsPeriodClosed(ctx context.Context, orgID uint, month, year int) (bool, error)
	ClosePeriod(ctx context.Context, period *model.ClosedPeriod) error
}

type closingRepository struct {
	db *gorm.DB
}

func NewClosingRepository(db *gorm.DB) ClosingRepository {
	return &closingRepository{db: db}
}

func (r *closingRepository) CreateLog(ctx context.Context, log *model.ClosingLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *closingRepository) UpdateLog(ctx context.Context, log *model.ClosingLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *closingRepository) GetLogByDate(ctx context.Context, orgID uint, date string, typ string) (*model.ClosingLog, error) {
	var log model.ClosingLog
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND date = ? AND type = ?", orgID, date, typ).
		First(&log).Error
	return &log, err
}

func (r *closingRepository) IsPeriodClosed(ctx context.Context, orgID uint, month, year int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ClosedPeriod{}).
		Where("organization_id = ? AND month = ? AND year = ?", orgID, month, year).
		Count(&count).Error
	return count > 0, err
}

func (r *closingRepository) ClosePeriod(ctx context.Context, period *model.ClosedPeriod) error {
	return r.db.WithContext(ctx).Create(period).Error
}
