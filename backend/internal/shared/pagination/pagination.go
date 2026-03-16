package pagination

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Params holds parsed pagination parameters.
type Params struct {
	Page  int
	Limit int
}

// Parse extracts page and limit from query string.
func Parse(c *fiber.Ctx) Params {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return Params{Page: page, Limit: limit}
}

// Offset returns the SQL offset for the current page.
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// TotalPages calculates the total number of pages.
func (p Params) TotalPages(totalItems int64) int {
	return int(math.Ceil(float64(totalItems) / float64(p.Limit)))
}

// Scope returns a GORM scope that applies limit and offset.
func (p Params) Scope() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset()).Limit(p.Limit)
	}
}
