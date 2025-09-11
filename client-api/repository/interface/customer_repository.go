package repository

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Clientes struct {
	ID    int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Nome  string          `gorm:"not null" json:"nome"`
	Email string          `gorm:"not null;unique" json:"email"`
	Saldo decimal.Decimal `gorm:"type:TEXT;not null" json:"saldo"`
}

type CustomerRepo interface {
	Migrate() error
	GetAll() ([]Clientes, error)
	GetByID(id int) (Clientes, error)
	Create(c *Clientes) error
	Update(c *Clientes) error
	Delete(id int) error
	ChangeBalance(id int, delta decimal.Decimal) (Clientes, error)
	DB() *gorm.DB
}
