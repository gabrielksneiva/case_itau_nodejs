package customer

import (
	repository "case-itau/repository/interface"
	"errors"
	"sync"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Service contains business logic for customers
type Service struct {
	db     *gorm.DB
	mutexs sync.Map
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) getLock(id uint) *sync.Mutex {
	v, ok := s.mutexs.Load(id)
	if ok {
		return v.(*sync.Mutex)
	}
	m := &sync.Mutex{}
	actual, _ := s.mutexs.LoadOrStore(id, m)
	return actual.(*sync.Mutex)
}

func (s *Service) ListAll() ([]repository.Clientes, error) {
	var list []repository.Clientes
	if err := s.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Service) GetByID(id uint) (repository.Clientes, error) {
	var c repository.Clientes
	if err := s.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		return repository.Clientes{}, err
	}
	return c, nil
}

func (s *Service) Create(input repository.Clientes) (repository.Clientes, error) {
	input.Saldo = decimal.Zero
	if err := s.db.Create(&input).Error; err != nil {
		return repository.Clientes{}, err
	}
	return input, nil
}

func (s *Service) Update(id uint, input repository.Clientes) (repository.Clientes, error) {
	var c repository.Clientes
	if err := s.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		return repository.Clientes{}, err
	}
	c.Nome = input.Nome
	c.Email = input.Email
	if err := s.db.Save(&c).Error; err != nil {
		return repository.Clientes{}, err
	}
	return c, nil
}

func (s *Service) Delete(id uint) error {
	if err := s.db.Delete(&repository.Clientes{}, id).Error; err != nil {
		return err
	}
	return nil
}

var ErrNotFound = errors.New("customer not found")

func (s *Service) Deposit(id uint, amount decimal.Decimal) (repository.Clientes, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return repository.Clientes{}, errors.New("amount must be greater than zero")
	}

	m := s.getLock(id)
	m.Lock()
	defer m.Unlock()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var c repository.Clientes
		if err := tx.Clauses().First(&c, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		c.Saldo = c.Saldo.Add(amount)
		if err := tx.Save(&c).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return repository.Clientes{}, err
	}

	return s.GetByID(id)
}

func (s *Service) Withdraw(id uint, amount decimal.Decimal) (repository.Clientes, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return repository.Clientes{}, errors.New("amount must be greater than zero")
	}

	m := s.getLock(id)
	m.Lock()
	defer m.Unlock()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var c repository.Clientes
		if err := tx.First(&c, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if c.Saldo.LessThan(amount) {
			return errors.New("insufficient balance")
		}
		c.Saldo = c.Saldo.Sub(amount)
		if err := tx.Save(&c).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return repository.Clientes{}, err
	}
	return s.GetByID(id)
}
