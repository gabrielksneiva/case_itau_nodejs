package customer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"case-itau/repositories"

	"github.com/shopspring/decimal"
)

var (
	ErrNotFound          = errors.New("cliente não encontrado")
	ErrInsufficientFunds = errors.New("saldo insuficiente")
	ErrUniqueEmail       = errors.New("UNIQUE constraint failed: clientes.email")
)

type Service struct {
	repo repositories.IRepository[repositories.Clientes]
}

func NewService(repo repositories.IRepository[repositories.Clientes]) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAll(ctx context.Context) ([]repositories.Clientes, error) {
	return s.repo.Find(ctx, nil, "", 0, 0)
}

func (s *Service) GetByID(ctx context.Context, id int) (*repositories.Clientes, error) {
	c, err := s.repo.FindOne(ctx, fmt.Sprintf("id = %d", id))
	if err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) Create(ctx context.Context, input repositories.Clientes) (repositories.Clientes, error) {
	input.Saldo = decimal.Zero

	if err := s.repo.InsertOne(ctx, &input); err != nil {
		if strings.Contains(err.Error(), ErrUniqueEmail.Error()) {
			return repositories.Clientes{}, ErrUniqueEmail
		}
		return repositories.Clientes{}, err
	}
	return input, nil
}

func (s *Service) Update(ctx context.Context, id int, input repositories.Clientes) (*repositories.Clientes, error) {
	update := make(map[string]any)
	if input.Nome != "" {
		update["nome"] = input.Nome
	}

	if input.Email != "" {
		update["email"] = input.Email
	}

	if err := s.repo.UpdateOne(ctx, fmt.Sprintf("id = %d", id), update); err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, err
		}
		if strings.Contains(err.Error(), ErrUniqueEmail.Error()) {
			return nil, ErrUniqueEmail
		}
		return nil, err
	}

	updated, err := s.repo.FindOne(ctx, fmt.Sprintf("id = %d", id))
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if err := s.repo.DeleteOne(ctx, id); err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Transactions(ctx context.Context, id int, delta decimal.Decimal) (*repositories.Clientes, error) {
	c, err := s.repo.FindOne(ctx, fmt.Sprintf("id = %d", id))
	if err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	newSaldo := c.Saldo.Add(delta)
	if newSaldo.IsNegative() {
		return nil, ErrInsufficientFunds
	}

	updates := map[string]any{
		"saldo": newSaldo,
	}

	err = s.repo.UpdateOne(ctx, fmt.Sprintf("id = %d", c.ID), updates)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.FindOne(ctx, fmt.Sprintf("id = %d", c.ID))
	if err != nil {
		return nil, err
	}

	return updated, nil
}
