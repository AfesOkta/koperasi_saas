package dto

type SubscriptionPlanRequest struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Code        string  `json:"code" validate:"required,max=50"`
	Description string  `json:"description" validate:"omitempty"`
	Price       float64 `json:"price" validate:"required,min=0"`
	MaxUsers    int     `json:"max_users" validate:"min=0"`
	MaxMembers  int     `json:"max_members" validate:"min=0"`
}

type SubscriptionPlanResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	MaxUsers    int     `json:"max_users"`
	MaxMembers  int     `json:"max_members"`
}

type OrgSubscriptionResponse struct {
	ID             uint                     `json:"id"`
	OrganizationID uint                     `json:"organization_id"`
	PlanID         uint                     `json:"plan_id"`
	Plan           SubscriptionPlanResponse `json:"plan"`
	StartDate      string                   `json:"start_date"`
	EndDate        string                   `json:"end_date"`
	Status         string                   `json:"status"`
}
