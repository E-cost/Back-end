package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePhone(fl validator.FieldLevel) bool {
	phoneRegex := `^\+\d{1,3}\d{9,14}$`
	phone := fl.Field().String()
	match, _ := regexp.MatchString(phoneRegex, phone)
	return match
}
