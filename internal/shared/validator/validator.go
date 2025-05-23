package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// validateSlug checks if the string is a valid slug.
func validateSlug(fl validator.FieldLevel) bool {
	return slugRegex.MatchString(fl.Field().String())
}

// validatePermissionCasbin checks if the permission is valid for Casbin.
func validatePermissionCasbin(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Check if the field is a slice
	if field.Kind() != reflect.Slice {
		return false
	}

	for i := range field.Len() {
		arrayElement := field.Index(i)

		// Check if the element is a slice
		if arrayElement.Kind() != reflect.Slice {
			return false
		}

		// Check length of the inner array
		if arrayElement.Len() != 2 {
			return false
		}
	}

	return true
}

// CustomValidator implements validation for Fiber
type CustomValidator struct {
	validator *validator.Validate
}

// New creates a new CustomValidator instance.
func New() *CustomValidator {
	v := validator.New()

	// Register custom tag name function to use json tags in error messages.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	// Register custom validations
	if err := v.RegisterValidation("slug", validateSlug); err != nil {
		panic("failed to register slug validation: " + err.Error())
	}
	if err := v.RegisterValidation("casbin_permission", validatePermissionCasbin); err != nil {
		panic("failed to register slug validation: " + err.Error())
	}

	return &CustomValidator{validator: v}
}

// ValidateStruct validates a struct using validator.v10 tags
func (cv *CustomValidator) Struct(s any) error {
	if err := cv.validator.Struct(s); err != nil {
		return err
	}
	return nil
}

// ValidateVar validates a single variable with a single rule
func (cv *CustomValidator) ValidateVar(field any, tag string) error {
	return cv.validator.Var(field, tag)
}

// Engine returns the underlying validator engine
func (cv *CustomValidator) Engine() *validator.Validate {
	return cv.validator
}

// RegisterValidation registers a custom validation function
func (cv *CustomValidator) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error {
	return cv.validator.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}
