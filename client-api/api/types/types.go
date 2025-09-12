package types

import (
	validations "case-itau/utils/validation"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
)

type CustomerDto struct {
	ID    int             `json:"id"`
	Nome  string          `json:"nome"`
	Email string          `json:"email"`
	Saldo decimal.Decimal `json:"saldo"`
}

type CreateCustomerRequest struct {
	Nome  string `json:"nome" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateCustomerRequest struct {
	Nome  string `json:"nome" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type TransactionRequest struct {
	Valor decimal.Decimal `json:"valor" validate:"required"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (fi *CreateCustomerRequest) IsValid(c *CreateCustomerRequest) error {
	return validations.Validate(c)
}

func (fi *CreateCustomerRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(&fi)
}

func (fi *UpdateCustomerRequest) IsValid(c *UpdateCustomerRequest) error {
	return validations.Validate(c)
}

func (fi *UpdateCustomerRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(&fi)
}

func (fi *TransactionRequest) IsValid(t *TransactionRequest) error {
	if t.Valor.LessThanOrEqual(decimal.Zero) {
		return errors.New("valor da transação deve ser maior que zero")
	}
	return validations.Validate(t)
}

func (fi *TransactionRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(fi)
}
