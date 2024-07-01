package dto

import (
	"Ecost/internal/utils/validation"
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CreateContactDto struct {
	Name         string `json:"name" validate:"required"`
	Surname      string `json:"surname"`
	Email        string `json:"email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required,phone"`
	Message      string `json:"message" validate:"dangerous"`
}

func (c *CreateContactDto) Validate() error {
	validate := validator.New()

	if err := validate.RegisterValidation("phone", validation.ValidatePhone); err != nil {
		return err
	}

	if err := validate.RegisterValidation("dangerous", validateDangerous); err != nil {
		return err
	}

	return validate.Struct(c)
}

func validateDangerous(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	dangerousPattern := `['";--]`
	match, _ := regexp.MatchString(dangerousPattern, value)

	return !match
}
