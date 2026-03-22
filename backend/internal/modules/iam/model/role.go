package model

// Role represents an RBAC role within an organization.
type Role struct {
	TenantModel
	Name        string       `json:"name" gorm:"not null;uniqueIndex:idx_roles_org_name"`
	Description string       `json:"description"`
	IsSystem    bool         `json:"is_system" gorm:"default:false"`
	Version     int          `json:"version" gorm:"not null;default:1"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
}
