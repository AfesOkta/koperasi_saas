package model

import "time"

// User represents an admin or system user within an organization.
type User struct {
	TenantModel
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"not null;index"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Phone        string    `json:"phone"`
	Avatar       string    `json:"avatar"`
	Status       string    `json:"status" gorm:"default:'active'"`
	LastLoginAt  time.Time `json:"last_login_at"`
	Roles        []Role    `json:"roles" gorm:"many2many:user_roles;"`
}
