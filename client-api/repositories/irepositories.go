package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrRepoNotFound         = errors.New("cliente não encontrado")
	ErrRepoInsufficientFund = errors.New("saldo insuficiente")
)

type Customers struct {
	ID      uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Name    string          `gorm:"not null" json:"name"`
	Email   string          `gorm:"not null;unique" json:"email"`
	Balance decimal.Decimal `gorm:"type:TEXT;not null" json:"balance"`
}

type Transaction struct {
	TransactionID uuid.UUID       `gorm:"type:uuid;primaryKey" json:"transaction_id"`
	CustomerID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"customer_id"`
	Customer      Customers       `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Amount        decimal.Decimal `gorm:"type:text;not null" json:"amount"`
	Type          string          `gorm:"type:text;not null" json:"type"`
	CreatedAt     time.Time       `gorm:"autoCreateTime" json:"created_at"`
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
