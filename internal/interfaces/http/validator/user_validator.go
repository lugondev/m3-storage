package validator

import (
	"github.com/go-playground/validator/v10"
)

func RegisterUserValidators(v *validator.Validate) {
	// Custom validation for user roles
	v.RegisterValidation("user_role", validateUserRole)
}

func validateUserRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == "admin" || role == "user"
}
