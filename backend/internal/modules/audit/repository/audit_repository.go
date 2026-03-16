package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/audit/model"
	"gorm.io/gorm"
)

type AuditRepository interface {
	CreateLog(ctx context.Context, log *model.AuditLog) error
	ListLogs(ctx context.Context, orgID uint, limit, offset int) ([]model.AuditLog, int64, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) CreateLog(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *auditRepository) ListLogs(ctx context.Context, orgID uint, limit, offset int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{}).Where("organization_id = ?", orgID)
	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error

	return logs, total, err
}
