package types

import (
	validations "case-itau/utils/validation"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CustomerDto struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Balance decimal.Decimal `json:"balance"`
}

type TransactionDto struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	CustomerID    uuid.UUID       `json:"customer_id"`
	Amount        decimal.Decimal `json:"value"`
	Type          string          `json:"type"`
	CreatedAt     time.Time       `json:"created_at"`
}

type CreateCustomerRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

type TransactionRequest struct {
	Amount decimal.Decimal `json:"amount" validate:"required"`
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
	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("valor da transação deve ser maior que zero")
	}
	return validations.Validate(t)
}

func (fi *TransactionRequest) FromBody(ctx *fiber.Ctx) error {
	return ctx.BodyParser(fi)
}
