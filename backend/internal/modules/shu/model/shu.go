package model

import (
	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

type SHUConfig struct {
	model.TenantModel
	Year              int     `json:"year" gorm:"not null;index"`
	TotalSHU          float64 `json:"total_shu" gorm:"type:decimal(15,2);not null"`
	MemberSavingsPct  float64 `json:"member_savings_pct" gorm:"type:decimal(5,2);not null"`  // Percentage for savings share
	MemberBusinessPct float64 `json:"member_business_pct" gorm:"type:decimal(5,2);not null"` // Percentage for business/transaction share
	Status            string  `json:"status" gorm:"size:20;default:'draft'"`                 // draft, calculated, distributed
}

type SHUDistribution struct {
	model.TenantModel
	SHUConfigID   uint    `json:"shu_config_id" gorm:"not null;index"`
	MemberID      uint    `json:"member_id" gorm:"not null;index"`
	SavingsShare  float64 `json:"savings_share" gorm:"type:decimal(15,2);not null"`
	BusinessShare float64 `json:"business_share" gorm:"type:decimal(15,2);not null"`
	TotalAmount   float64 `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Status        string  `json:"status" gorm:"size:20;default:'pending'"` // pending, paid
}
