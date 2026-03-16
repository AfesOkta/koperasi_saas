package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// Tenant middleware ensures the organization_id is present in the context.
// Must be used after Auth middleware.
func Tenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := GetOrganizationID(c)
		if orgID == 0 {
			return response.Forbidden(c, "Organization context not found")
		}
		return c.Next()
	}
}
