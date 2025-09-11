package customer

import (
	"errors"

	repository "case-itau/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrNotFound          = errors.New("cliente não encontrado")
	ErrInsufficientFunds = errors.New("saldo insuficiente")
	ErrConflict          = errors.New("conflito de concorrência")
)

type Service struct {
	repo repository.CustomerRepo
}

func NewService(repo repository.CustomerRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAll() ([]repository.Clientes, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (repository.Clientes, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		return repository.Clientes{}, err
	}
	return c, nil
}

func (s *Service) Create(input repository.Clientes) (repository.Clientes, error) {
	input.Saldo = decimal.Zero
	if input.Version == 0 {
		input.Version = 1
	}
	if err := s.repo.Create(&input); err != nil {
		return repository.Clientes{}, err
	}
	return input, nil
}

func (s *Service) Update(id int, input repository.Clientes) (repository.Clientes, error) {
	curr, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		return repository.Clientes{}, err
	}

	curr.Nome = input.Nome
	curr.Email = input.Email

	if err := s.repo.Update(&curr); err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		if errors.Is(err, repository.ErrRepoConflict) {
			return repository.Clientes{}, ErrConflict
		}
		return repository.Clientes{}, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return repository.Clientes{}, err
	}
	return updated, nil
}

func (s *Service) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrRepoNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Service) ChangeBalanceWithPessimisticLock(id int, delta decimal.Decimal) (repository.Clientes, error) {
	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return repository.Clientes{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var c repository.Clientes
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.Clientes{}, ErrNotFound
		}
		return repository.Clientes{}, err
	}

	newSaldo := c.Saldo.Add(delta)
	if newSaldo.IsNegative() {
		tx.Rollback()
		return repository.Clientes{}, ErrInsufficientFunds
	}
	c.Saldo = newSaldo

	if err := tx.Save(&c).Error; err != nil {
		tx.Rollback()
		return repository.Clientes{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return repository.Clientes{}, err
	}

	return c, nil
}

func (s *Service) Deposit(id int, amount decimal.Decimal) (repository.Clientes, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return repository.Clientes{}, errors.New("amount deve ser maior que zero")
	}
	return s.ChangeBalanceWithPessimisticLock(id, amount)
}

func (s *Service) Withdraw(id int, amount decimal.Decimal) (repository.Clientes, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return repository.Clientes{}, errors.New("amount deve ser maior que zero")
	}
	return s.ChangeBalanceWithPessimisticLock(id, amount.Neg())
}
