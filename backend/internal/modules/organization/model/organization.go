package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Organization represents a tenant in the multi-tenant SaaS.
type Organization struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Logo      string         `json:"logo"`
	Plan      string         `json:"plan" gorm:"default:'basic'"`
	Settings  datatypes.JSON `json:"settings" gorm:"type:jsonb;default:'{}'"`
	Status    string         `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
