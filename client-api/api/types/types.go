package types

import (
	validations "case-itau/utils/validation"
	"errors"

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
	Amount decimal.Decimal `json:"amount" validate:"required"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func IsValidCreateCustomerRequest(c CreateCustomerRequest) error {
	if err := validations.Validate(c); err != nil {
		return err
	}
	return nil
}

func IsValidUpdateCustomerRequest(c UpdateCustomerRequest) error {
	if err := validations.Validate(c); err != nil {
		return err
	}
	return nil
}

func IsValidTransactionRequest(t TransactionRequest) error {
	if err := validations.Validate(t); err != nil {
		return err
	}

	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("valor da transação deve ser maior que zero")
	}
	return nil
}
