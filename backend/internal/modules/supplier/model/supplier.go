package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// Supplier represents an entity that the cooperative purchases inventory from.
type Supplier struct {
	model.TenantModel
	Code        string `json:"code" gorm:"not null;uniqueIndex:idx_supplier_code_org"`
	Name        string `json:"name" gorm:"not null"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Status      string `json:"status" gorm:"default:'active'"` // active, inactive
}
