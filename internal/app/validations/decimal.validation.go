package validations

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateDecimal(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	regex := `^\d{4}(\.\d{1,10})?$`
	match, _ := regexp.MatchString(regex, value)
	return match
}
