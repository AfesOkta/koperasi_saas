package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/modules/closing/repository"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// PeriodGuard prevents creating or modifying transactions in closed accounting periods.
func PeriodGuard(repo repository.ClosingRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only check for write operations
		method := strings.ToUpper(c.Method())
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return c.Next()
		}

		// Skip for SuperAdmins (optional, based on decision 4)
		role := c.Locals("role").(string)
		if role == "SUPER_ADMIN" {
			return c.Next()
		}

		orgID := GetOrganizationID(c)
		
		// For transactions, we normally check the transaction date.
		// If the request doesn't provide a date, we assume today.
		// In a real implementation, we'd parse the request body for "date".
		// For MVP, we'll check the current month if no date is specified.
		
		now := time.Now()
		month := int(now.Month())
		year := now.Year()

		// Attempt to grab date from query or body if helpful
		// (Advanced: can use a middleware that pre-parses specific fields)

		isClosed, err := repo.IsPeriodClosed(c.Context(), orgID, month, year)
		if err != nil {
			return response.InternalError(c, "Failed to verify period status")
		}

		if isClosed {
			return response.Forbidden(c, "This accounting period is already closed. Contact a Super Admin for overrides.")
		}

		return c.Next()
	}
}
