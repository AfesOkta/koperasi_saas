package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// RBAC returns a middleware that checks if the user has the required permission.
// Permissions are checked against a permission checker function that should
// be implemented by the IAM module.
type PermissionChecker func(roleID uint, permission string) bool

// RequirePermission checks if the authenticated user has a specific permission.
func RequirePermission(permission string, checker PermissionChecker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("role_id").(uint)
		if !ok || roleID == 0 {
			return response.Forbidden(c, "No role assigned")
		}

		if !checker(roleID, permission) {
			return response.Forbidden(c, "Insufficient permissions")
		}

		return c.Next()
	}
}
