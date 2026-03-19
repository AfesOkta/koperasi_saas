package dto

import "time"

type SHUConfigRequest struct {
	Year              int     `json:"year" validate:"required"`
	TotalSHU          float64 `json:"total_shu" validate:"required,gt=0"`
	MemberSavingsPct  float64 `json:"member_savings_pct" validate:"required,min=0,max=100"`
	MemberBusinessPct float64 `json:"member_business_pct" validate:"required,min=0,max=100"`
}

type SHUConfigResponse struct {
	ID                uint      `json:"id"`
	Year              int       `json:"year"`
	TotalSHU          float64   `json:"total_shu"`
	MemberSavingsPct  float64   `json:"member_savings_pct"`
	MemberBusinessPct float64   `json:"member_business_pct"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type SHUDistributionResponse struct {
	ID            uint      `json:"id"`
	MemberID      uint      `json:"member_id"`
	SavingsShare  float64   `json:"savings_share"`
	BusinessShare float64   `json:"business_share"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
