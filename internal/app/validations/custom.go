package validations

import "github.com/go-playground/validator/v10"

func RegisterCustomValidation(v *validator.Validate) {
	v.RegisterValidation("decimal", ValidateDecimal)
}
