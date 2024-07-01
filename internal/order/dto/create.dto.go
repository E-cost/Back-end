package dto

import (
	"Ecost/internal/utils/validation"

	"github.com/go-playground/validator/v10"
)

type CreateOrderDto struct {
	Name         string `json:"name" validate:"required"`
	Surname      string `json:"surname" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	ContactPhone string `json:"contact_phone" validate:"required,phone"`
	Items        []Item `json:"items" validate:"required"`
}

type Item struct {
	ID       string `json:"id" validate:"required"`
	Category string `json:"category" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

func (c *CreateOrderDto) Validate() error {
	validate := validator.New()
	if err := validate.RegisterValidation("phone", validation.ValidatePhone); err != nil {
		return err
	}
	return validate.Struct(c)
}
