package model

import (
	"time"

	"github.com/koperasi-gresik/backend/internal/modules/iam/model"
)

// ClosingLog tracks the execution of EOD and EOM processes.
type ClosingLog struct {
	model.TenantModel
	Date         string     `json:"date" gorm:"type:date;not null;index"` // YYYY-MM-DD
	Type         string     `json:"type" gorm:"size:20;not null;index"`   // EOD, EOM
	Status       string     `json:"status" gorm:"size:20;default:'PENDING'"` // PENDING, RUNNING, SUCCESS, FAILED
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	ErrorMessage string     `json:"error_message" gorm:"type:text"`
	ProcessedBy  uint       `json:"processed_by" gorm:"index"` // 0 for system/cron
}

// ClosedPeriod prevents transactions in closed accounting months.
type ClosedPeriod struct {
	model.TenantModel
	Month    int    `json:"month" gorm:"not null"`
	Year     int    `json:"year" gorm:"not null"`
	ClosedAt string `json:"closed_at" gorm:"type:timestamp;not null"`
	ClosedBy uint   `json:"closed_by" gorm:"not null"`
}
