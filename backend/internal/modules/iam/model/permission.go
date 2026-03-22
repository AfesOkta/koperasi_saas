package model

import "time"

// Permission represents an RBAC permission in the system.
// System-wide, un-scoped by organization.
// Name format: resource:action or resource:action:scope
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex"`
	Resource    string    `json:"resource" gorm:"not null"`
	Action      string    `json:"action" gorm:"not null"`
	Scope       string    `json:"scope" gorm:"not null;default:'any'"` // 'any' or 'own'
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
