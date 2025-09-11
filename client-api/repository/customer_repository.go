package repository

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrRepoNotFound         = errors.New("cliente não encontrado")
	ErrRepoConflict         = errors.New("conflito de concorrência")
	ErrRepoInsufficientFund = errors.New("saldo insuficiente")
)

type Clientes struct {
	ID      int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome    string          `gorm:"not null" json:"nome"`
	Email   string          `gorm:"not null;unique" json:"email"`
	Saldo   decimal.Decimal `gorm:"type:TEXT;not null" json:"saldo"`
	Version int             `gorm:"not null;default:1" json:"version"`
}

type CustomerRepo interface {
	Migrate() error
	Create(c *Clientes) error
	GetAll() ([]Clientes, error)
	GetByID(id int) (Clientes, error)
	Update(c *Clientes) error
	Delete(id int) error
	DB() *gorm.DB
}

var _ CustomerRepo = (*GormRepository[Clientes, int])(nil)

func NewCustomerRepository(db *gorm.DB) CustomerRepo {
	return NewGormRepository[Clientes, int](db)
}
