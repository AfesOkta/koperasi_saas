package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/koperasi-gresik/backend/internal/shared/response"
)

var validate = validator.New()

// Validate validates a struct and returns a Fiber error response if invalid.
func Validate(c *fiber.Ctx, s interface{}) error {
	if err := c.BodyParser(s); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := validate.Struct(s); err != nil {
		errors := formatErrors(err.(validator.ValidationErrors))
		return response.BadRequest(c, "Validation failed", errors)
	}

	return nil
}

// ValidateStruct validates a struct without parsing (for non-HTTP use).
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func formatErrors(errs validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for _, err := range errs {
		result[err.Field()] = formatMessage(err)
	}
	return result
}

func formatMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return err.Field() + " must be at least " + err.Param()
	case "max":
		return err.Field() + " must be at most " + err.Param()
	case "gte":
		return err.Field() + " must be greater than or equal to " + err.Param()
	case "lte":
		return err.Field() + " must be less than or equal to " + err.Param()
	default:
		return err.Field() + " failed " + err.Tag() + " validation"
	}
}
