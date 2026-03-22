package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// RequirePermission returns a middleware that checks if the authenticated user's role
// has the given permission using the Redis-backed PermissionCache.
// Flow: JWT role_id + role_version → Redis key → in-memory check (zero DB queries).
func RequirePermission(permission string, cache *PermissionCache) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID, ok := c.Locals("role_id").(uint)
		if !ok || roleID == 0 {
			return response.Forbidden(c, "No role assigned")
		}

		roleVersion, _ := c.Locals("role_version").(int)

		if !cache.HasPermission(c.Context(), roleID, roleVersion, permission) {
			// Log the denial (audit integration point)
			// This will be wired to audit logger in a subsequent step
			return response.Forbidden(c, "Permission denied: "+permission)
		}

		return c.Next()
	}
}

// GetRoleVersion extracts role_version from context.
func GetRoleVersion(c *fiber.Ctx) int {
	if v, ok := c.Locals("role_version").(int); ok {
		return v
	}
	return 0
}
