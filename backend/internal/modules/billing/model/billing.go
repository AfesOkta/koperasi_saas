package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

type SubscriptionPlan struct {
	model.BaseModel
	Name        string  `json:"name" gorm:"size:100;not null;unique"`
	Code        string  `json:"code" gorm:"size:50;not null;unique"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"type:decimal(15,2);not null"`
	MaxUsers    int     `json:"max_users" gorm:"default:0"` // 0 = unlimited
	MaxMembers  int     `json:"max_members" gorm:"default:0"`
}

type OrgSubscription struct {
	model.BaseModel
	OrganizationID uint             `json:"organization_id" gorm:"not null;uniqueIndex"`
	PlanID         uint             `json:"plan_id" gorm:"not null"`
	Plan           SubscriptionPlan `json:"plan" gorm:"foreignKey:PlanID"`
	StartDate      time.Time        `json:"start_date" gorm:"not null"`
	EndDate        time.Time        `json:"end_date" gorm:"not null"`
	Status         string           `json:"status" gorm:"size:20;default:'active'"` // active, expired, cancelled
	RenewalDate    *time.Time       `json:"renewal_date"`
}
