package model

import "github.com/koperasi-gresik/backend/internal/modules/iam/model"

// Member represents a cooperative member.
type Member struct {
	model.TenantModel
	UserID       uint             `json:"user_id" gorm:"index"` // Link to IAM login
	MemberNumber string           `json:"member_number" gorm:"not null;uniqueIndex:idx_member_no_org"`
	Name         string           `json:"name" gorm:"not null"`
	NIK          string           `json:"nik" gorm:"not null;uniqueIndex:idx_member_nik_org"`
	Address      string           `json:"address"`
	Phone        string           `json:"phone"`
	Status       string           `json:"status" gorm:"default:'pending'"` // pending, active, inactive
	Documents    []MemberDocument `json:"documents,omitempty" gorm:"foreignKey:MemberID"`
	Cards        []MemberCard     `json:"cards,omitempty" gorm:"foreignKey:MemberID"`
}

// MemberDocument represents uploaded KYC documents.
type MemberDocument struct {
	model.TenantModel
	MemberID uint   `json:"member_id" gorm:"not null;index"`
	Type     string `json:"type" gorm:"not null"` // e.g., ktp, kk, selfie
	FileURL  string `json:"file_url" gorm:"not null"`
}

// MemberCard represents the physical/digital ID card.
type MemberCard struct {
	model.TenantModel
	MemberID   uint   `json:"member_id" gorm:"not null;index"`
	CardNumber string `json:"card_number" gorm:"not null;uniqueIndex"`
	Status     string `json:"status" gorm:"default:'active'"` // active, blocked, lost
}
