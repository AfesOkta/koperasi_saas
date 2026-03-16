package dto

// OrganizationCreateRequest DTO for creating an organization.
type OrganizationCreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"omitempty,max=20"`
	Address string `json:"address"`
}

// OrganizationUpdateRequest DTO for updating an organization.
type OrganizationUpdateRequest struct {
	Name     string                 `json:"name" validate:"omitempty,min=2,max=255"`
	Email    string                 `json:"email" validate:"omitempty,email"`
	Phone    string                 `json:"phone" validate:"omitempty,max=20"`
	Address  string                 `json:"address"`
	Plan     string                 `json:"plan" validate:"omitempty,oneof=free basic premium"`
	Settings map[string]interface{} `json:"settings"`
}

// OrganizationResponse DTO for returning organization data.
type OrganizationResponse struct {
	ID        uint                   `json:"id"`
	Name      string                 `json:"name"`
	Slug      string                 `json:"slug"`
	Email     string                 `json:"email"`
	Phone     string                 `json:"phone"`
	Address   string                 `json:"address"`
	Logo      string                 `json:"logo"`
	Plan      string                 `json:"plan"`
	Settings  map[string]interface{} `json:"settings"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"created_at"`
}

// OnboardingRequest DTO for registration of a new cooperative with its first admin.
type OnboardingRequest struct {
	// Organization Details
	OrganizationName string `json:"organization_name" validate:"required,min=2,max=255"`
	Email            string `json:"email" validate:"required,email"`
	Phone            string `json:"phone" validate:"omitempty,max=20"`
	Address          string `json:"address"`

	// Admin User Details
	AdminName     string `json:"admin_name" validate:"required,min=2,max=255"`
	AdminEmail    string `json:"admin_email" validate:"required,email"`
	AdminPassword string `json:"admin_password" validate:"required,min=6"`
}

// OnboardingResponse DTO for returning onboarding result.
type OnboardingResponse struct {
	Organization OrganizationResponse `json:"organization"`
	AdminUser    interface{}          `json:"admin_user"` // Using interface{} to avoid circular dependency if needed or use a light struct
}
