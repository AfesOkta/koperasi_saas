package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// SuperAdminGuard ensures the request is from Organization ID 1 (Platform Owner).
func SuperAdminGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := GetOrganizationID(c)
		if orgID != 1 {
			return response.Forbidden(c, "SuperAdmin access required")
		}
		return c.Next()
	}
}
