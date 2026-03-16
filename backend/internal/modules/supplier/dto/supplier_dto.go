package dto

type SupplierCreateRequest struct {
	Code        string `json:"code" validate:"required"`
	Name        string `json:"name" validate:"required"`
	ContactName string `json:"contact_name" validate:"omitempty"`
	Phone       string `json:"phone" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty,email"`
	Address     string `json:"address" validate:"omitempty"`
}

type SupplierUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	ContactName string `json:"contact_name" validate:"omitempty"`
	Phone       string `json:"phone" validate:"omitempty"`
	Email       string `json:"email" validate:"omitempty,email"`
	Address     string `json:"address" validate:"omitempty"`
	Status      string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type SupplierResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}
