package repositories

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrRepoNotFound         = errors.New("cliente não encontrado")
	ErrRepoInsufficientFund = errors.New("saldo insuficiente")
)

type Clientes struct {
	ID      int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome    string          `gorm:"not null" json:"nome"`
	Email   string          `gorm:"not null;unique" json:"email"`
	Saldo   decimal.Decimal `gorm:"type:TEXT;not null" json:"saldo"`
}

type IRepository[T any] interface {
	WithPreload(associations ...string) *gormRepository[T]
	Find(ctx context.Context, where any, order string, limit, offset int) ([]T, error)
	FindOne(ctx context.Context, where any) (*T, error)
	InsertOne(ctx context.Context, entity *T) error
	UpdateOne(ctx context.Context, where any, updates map[string]any) error
	DeleteOne(ctx context.Context, where any) error
	Count(ctx context.Context, where any) (int64, error)
}
