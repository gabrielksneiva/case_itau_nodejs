package customer

import (
	"context"
	"errors"
	"strings"

	"case-itau/repositories"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound          = errors.New("cliente não encontrado")
	ErrInsufficientFunds = errors.New("saldo insuficiente")
	ErrUniqueEmail       = errors.New("UNIQUE constraint failed: customers.email")
)

type Service struct {
	repoCli   repositories.IRepository[repositories.Customers]
	repoTrans repositories.IRepository[repositories.Transaction]
}

func NewService(repoCli repositories.IRepository[repositories.Customers], repoTrans repositories.IRepository[repositories.Transaction]) *Service {
	return &Service{repoCli: repoCli, repoTrans: repoTrans}
}

func (s *Service) ListAll(ctx context.Context) ([]repositories.Customers, error) {
	return s.repoCli.Find(ctx, nil, "", 0, 0)
}

func (s *Service) GetByID(ctx context.Context, id string) (*repositories.Customers, error) {
	c, err := s.repoCli.FindOne(ctx, map[string]any{"id": id})
	if err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (s *Service) Create(ctx context.Context, input repositories.Customers) (repositories.Customers, error) {
	input.Balance = decimal.Zero

	if err := s.repoCli.InsertOne(ctx, &input); err != nil {
		if strings.Contains(err.Error(), ErrUniqueEmail.Error()) {
			return repositories.Customers{}, ErrUniqueEmail
		}
		return repositories.Customers{}, err
	}
	return input, nil
}

func (s *Service) Update(ctx context.Context, id string, input repositories.Customers) (*repositories.Customers, error) {
	update := make(map[string]any)
	if input.Name != "" {
		update["name"] = input.Name
	}

	if input.Email != "" {
		update["email"] = input.Email
	}

	if err := s.repoCli.UpdateOne(ctx, map[string]any{"id": id}, update); err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, err
		}
		if strings.Contains(err.Error(), ErrUniqueEmail.Error()) {
			return nil, ErrUniqueEmail
		}
		return nil, err
	}

	updated, err := s.repoCli.FindOne(ctx, map[string]any{"id": id})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repoCli.DeleteOne(ctx, map[string]any{"id": id}); err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Transactions(ctx context.Context, id string, delta decimal.Decimal) (*repositories.Customers, error) {
	c, err := s.repoCli.FindOne(ctx, map[string]any{"id": id})
	if err != nil {
		if errors.Is(err, repositories.ErrRepoNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	newBalance := c.Balance.Add(delta)
	if newBalance.IsNegative() {
		return nil, ErrInsufficientFunds
	}

	updates := map[string]any{
		"balance": newBalance,
	}

	err = s.repoCli.UpdateOne(ctx, map[string]any{"id": c.ID.String()}, updates)
	if err != nil {
		return nil, err
	}

	updated, err := s.repoCli.FindOne(ctx, map[string]any{"id": c.ID.String()})
	if err != nil {
		return nil, err
	}

	transactionType := "withdraw"
	if delta.IsPositive() {
		transactionType = "deposit"
	}

	t := &repositories.Transaction{
		TransactionID: uuid.New(),
		CustomerID:    c.ID,
		Amount:        delta,
		Type:          transactionType,
	}

	if err := s.repoTrans.InsertOne(ctx, t); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) ListTransactions(ctx context.Context, customerID string, page, size int) ([]repositories.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	offset := (page - 1) * size

	total, err := s.repoTrans.Count(ctx, map[string]any{"customer_id": customerID})
	if err != nil {
		return nil, 0, err
	}

	txs, err := s.repoTrans.Find(
		ctx,
		map[string]any{"customer_id": customerID},
		"created_at DESC",
		size,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}
