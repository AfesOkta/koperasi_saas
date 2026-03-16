package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel is the base model for all domain entities.
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// TenantModel extends BaseModel with organization_id for multi-tenant isolation.
type TenantModel struct {
	BaseModel
	OrganizationID uint `json:"organization_id" gorm:"index;not null"`
}
