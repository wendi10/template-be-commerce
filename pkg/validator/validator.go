package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct and returns a formatted error string.
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return formatErrors(err)
	}
	return nil
}

func formatErrors(err error) error {
	var errs []string
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, fieldError(e))
	}
	return fmt.Errorf("%s", strings.Join(errs, "; "))
}

func fieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, e.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
