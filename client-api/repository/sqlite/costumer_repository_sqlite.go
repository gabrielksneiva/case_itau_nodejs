package sqlite

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	repository "case-itau/repository/interface"
)

var _ repository.CustomerRepo = (*sqliteRepo)(nil)

// sqliteRepo implements CustomerRepo
type sqliteRepo struct {
	db *gorm.DB
}

// Cria novo repositório SQLite implementando CustomerRepo
func NewSqliteCustomerRepo(db *gorm.DB) *sqliteRepo {
	return &sqliteRepo{db: db}
}

// DB retorna a conexão
func (r *sqliteRepo) DB() *gorm.DB {
	return r.db
}

// Migrate cria a tabela no DB
func (r *sqliteRepo) Migrate() error {
	return r.db.AutoMigrate(&repository.Clientes{})
}

// GetAll retorna todos os clientes
func (r *sqliteRepo) GetAll() ([]repository.Clientes, error) {
	var customers []repository.Clientes
	if err := r.db.Find(&customers).Error; err != nil {
		return nil, err
	}

	out := make([]repository.Clientes, len(customers))
	for i, c := range customers {
		out[i] = repository.Clientes{
			ID:    c.ID,
			Nome:  c.Nome,
			Email: c.Email,
			Saldo: c.Saldo,
		}
	}
	return out, nil
}

// GetByID retorna cliente pelo ID
func (r *sqliteRepo) GetByID(id int) (repository.Clientes, error) {
	var c repository.Clientes
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.Clientes{}, errors.New("cliente não encontrado")
		}
		return repository.Clientes{}, err
	}
	return repository.Clientes{
		ID:    c.ID,
		Nome:  c.Nome,
		Email: c.Email,
		Saldo: c.Saldo,
	}, nil
}

// Create adiciona novo cliente
func (r *sqliteRepo) Create(c *repository.Clientes) error {
	newC := repository.Clientes{
		Nome:  c.Nome,
		Email: c.Email,
		Saldo: c.Saldo,
	}
	if err := r.db.Create(&newC).Error; err != nil {
		return err
	}
	c.ID = newC.ID
	return nil
}

// Update altera os dados do cliente
func (r *sqliteRepo) Update(c *repository.Clientes) error {
	var existing repository.Clientes
	if err := r.db.First(&existing, c.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cliente não encontrado")
		}
		return err
	}

	existing.Nome = c.Nome
	existing.Email = c.Email
	// Saldo não é atualizado aqui, usar ChangeBalance
	if err := r.db.Save(&existing).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove cliente pelo ID
func (r *sqliteRepo) Delete(id int) error {
	if err := r.db.Delete(&repository.Clientes{}, id).Error; err != nil {
		return err
	}
	return nil
}

// ChangeBalance altera o saldo do cliente de forma segura
func (r *sqliteRepo) ChangeBalance(id int, delta decimal.Decimal) (repository.Clientes, error) {
	var c repository.Clientes
	tx := r.db.Begin()
	if err := tx.First(&c, id).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repository.Clientes{}, errors.New("cliente não encontrado")
		}
		return repository.Clientes{}, err
	}

	c.Saldo = c.Saldo.Add(delta)
	if c.Saldo.IsNegative() {
		tx.Rollback()
		return repository.Clientes{}, errors.New("saldo insuficiente")
	}

	if err := tx.Save(&c).Error; err != nil {
		tx.Rollback()
		return repository.Clientes{}, err
	}
	tx.Commit()

	return repository.Clientes{
		ID:    c.ID,
		Nome:  c.Nome,
		Email: c.Email,
		Saldo: c.Saldo,
	}, nil
}
