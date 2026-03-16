package model

import (
	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

type AuditLog struct {
	model.TenantModel
	UserID     uint   `json:"user_id" gorm:"index"`
	Action     string `json:"action" gorm:"size:50;not null;index"`    // create, update, delete, login, etc.
	Resource   string `json:"resource" gorm:"size:100;not null;index"` // users, members, savings, etc.
	ResourceID string `json:"resource_id" gorm:"size:100;index"`
	OldValues  string `json:"old_values" gorm:"type:text"`
	NewValues  string `json:"new_values" gorm:"type:text"`
	IPAddress  string `json:"ip_address" gorm:"size:45"`
	UserAgent  string `json:"user_agent" gorm:"type:text"`
}
