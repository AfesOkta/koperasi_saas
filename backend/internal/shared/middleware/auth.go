package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

// Claims represents JWT token claims.
type Claims struct {
	UserID         uint   `json:"user_id"`
	OrganizationID uint   `json:"organization_id"`
	Email          string `json:"email"`
	RoleID         uint   `json:"role_id"`
	RoleVersion    int    `json:"role_version"` // Used for Redis cache key: perm:{role_id}:{role_version}
	jwt.RegisteredClaims
}

// Auth returns a middleware that validates JWT tokens.
func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization format")
		}

		tokenStr := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Store claims in context for downstream handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("organization_id", claims.OrganizationID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID)
		c.Locals("role_version", claims.RoleVersion)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// GetUserID extracts user ID from context.
func GetUserID(c *fiber.Ctx) uint {
	if id, ok := c.Locals("user_id").(uint); ok {
		return id
	}
	return 0
}

// GetOrganizationID extracts organization ID from context.
func GetOrganizationID(c *fiber.Ctx) uint {
	if id, ok := c.Locals("organization_id").(uint); ok {
		return id
	}
	return 0
}
