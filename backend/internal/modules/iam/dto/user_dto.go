package dto

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"omitempty"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Name   string `json:"name" validate:"omitempty"`
	Email  string `json:"email" validate:"omitempty,email"`
	Phone  string `json:"phone" validate:"omitempty"`
	Status string `json:"status" validate:"omitempty,oneof=active inactive"`
	RoleID uint   `json:"role_id" validate:"omitempty"`
}

type UserResponse struct {
	ID             uint           `json:"id"`
	OrganizationID uint           `json:"organization_id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	Avatar         string         `json:"avatar"`
	Status         string         `json:"status"`
	Roles          []RoleResponse `json:"roles,omitempty"`
}

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    bool   `json:"is_system"`
}
