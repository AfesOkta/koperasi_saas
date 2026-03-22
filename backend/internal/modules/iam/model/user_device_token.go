package model

import "time"

type UserDeviceToken struct {
	TenantModel
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	DeviceToken string    `json:"device_token" gorm:"not null;uniqueIndex"`
	DeviceType  string    `json:"device_type"` // ios, android, web
	LastSeenAt  time.Time `json:"last_seen_at"`
}
