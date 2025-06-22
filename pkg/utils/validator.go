package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

/* ---------------------------- Simple validation --------------------------- */

func MustValid(v any) {
	validator := validator.New()
	if err := validator.Struct(v); err != nil {
		panic(fmt.Errorf("validation failed: %w", err))
	}
}

func Validate(v any) error {
	validator := validator.New()
	if err := validator.Struct(v); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

/* ---------------------------- Struct validation --------------------------- */

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(v any) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
