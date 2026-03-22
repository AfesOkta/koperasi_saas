package dto

type DeviceTokenRequest struct {
	DeviceToken string `json:"device_token" validate:"required"`
	DeviceType  string `json:"device_type" validate:"required,oneof=ios android web"`
}
