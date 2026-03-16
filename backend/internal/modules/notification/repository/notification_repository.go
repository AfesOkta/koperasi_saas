package repository

import (
	"context"

	"github.com/koperasi-gresik/backend/internal/modules/notification/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	ListByUser(ctx context.Context, orgID, userID uint, limit, offset int) ([]model.Notification, int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *model.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) ListByUser(ctx context.Context, orgID, userID uint, limit, offset int) ([]model.Notification, int64, error) {
	var ns []model.Notification
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID)
	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ns).Error

	return ns, total, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}
