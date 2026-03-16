package model

import (
	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

type Notification struct {
	model.TenantModel
	UserID  uint   `json:"user_id" gorm:"not null;index"`
	Title   string `json:"title" gorm:"size:255;not null"`
	Message string `json:"message" gorm:"type:text;not null"`
	Type    string `json:"type" gorm:"size:50;default:'info'"` // info, success, warning, danger
	IsRead  bool   `json:"is_read" gorm:"default:false;index"`
	Link    string `json:"link" gorm:"size:255"` // Optional deep link
}
