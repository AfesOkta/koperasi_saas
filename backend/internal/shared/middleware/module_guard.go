package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
	"gorm.io/gorm"
)

// ModuleGuard middleware checks if a module is enabled for the organization.
func ModuleGuard(db *gorm.DB, moduleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orgID := GetOrganizationID(c)
		if orgID == 0 {
			return response.Forbidden(c, "Organization context not found")
		}

		// In a real application, you'd cache this or get it from a pre-loaded organization object
		var settings map[string]interface{}
		var settingsRaw []byte

		err := db.Table("organizations").Select("settings").Where("id = ?", orgID).Row().Scan(&settingsRaw)
		if err != nil {
			return response.InternalError(c, "Failed to check module status")
		}

		if len(settingsRaw) > 0 {
			json.Unmarshal(settingsRaw, &settings)
		}

		enabledModules, ok := settings["enabled_modules"].([]interface{})
		if !ok {
			// If not defined, assume some default or restricted set for MVP
			// For now, let's just return forbidden if no modules are explicitly enabled
			return response.Forbidden(c, "Module not enabled for this organization")
		}

		enabled := false
		for _, m := range enabledModules {
			if strings.EqualFold(m.(string), moduleName) {
				enabled = true
				break
			}
		}

		if !enabled {
			return response.Forbidden(c, "Module '"+moduleName+"' is not part of your subscription")
		}

		return c.Next()
	}
}
