package types

import "github.com/shopspring/decimal"

type CustomerDto struct {
	ID    int            `json:"id"`
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
	Message string `json:"message"`
}
