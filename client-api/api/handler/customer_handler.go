package handler

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"case-itau/api/types"
	repo "case-itau/repositories"
	"case-itau/services/customer"
)

type CustomerHandler struct {
	service *customer.Service
}

func NewCustomerHandler(s *customer.Service) *CustomerHandler {
	return &CustomerHandler{
		service: s,
	}
}

// GetCustomers godoc
// @Summary      Lista todos os usuários
// @Description  Endpoint para listar todos os usuários
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes [get]
func (h *CustomerHandler) List(c *fiber.Ctx) error {
	list, err := h.service.ListAll(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}

	out := make([]types.CustomerDto, 0, len(list))
	for _, it := range list {
		out = append(out, types.CustomerDto{
			ID: it.ID, Name: it.Name, Email: it.Email, Balance: it.Balance,
		})
	}
	return c.JSON(out)
}

// GetCustomer godoc
// @Summary      Obtém um usuário pelo ID
// @Description  Endpoint para obter um usuário pelo ID
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do usuário"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [get]
func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	cust, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Name: cust.Name, Email: cust.Email, Balance: cust.Balance}
	return c.JSON(out)
}

// CreateCustomer godoc
// @Summary      Cria um novo usuário
// @Description  Endpoint para criar um novo usuário
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        customer  body      types.CreateCustomerRequest  true  "Dados do usuário"
// @Success      201  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes [post]
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	req := &types.CreateCustomerRequest{}
	err := req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}

	if err = req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	model := repo.Customers{
		ID:      uuid.New(),
		Name:    req.Name,
		Email:   req.Email,
		Balance: decimal.Zero,
	}
	created, err := h.service.Create(c.UserContext(), model)
	if err != nil {
		if errors.Is(err, customer.ErrUniqueEmail) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "EMAIL_ALREADY_EXISTS", Message: "Email já registrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: created.ID, Name: created.Name, Email: created.Email, Balance: created.Balance}

	return c.Status(fiber.StatusCreated).JSON(out)
}

// UpdateCustomer godoc
// @Summary      Atualiza um usuário existente
// @Description  Endpoint para atualizar um usuário existente
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "ID do usuário"
// @Param        customer  body      types.UpdateCustomerRequest  true  "Dados do usuário"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [put]
func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	req := &types.CreateCustomerRequest{}
	err := req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	in := repo.Customers{Name: req.Name, Email: req.Email}
	updated, err := h.service.Update(c.UserContext(), id, in)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrUniqueEmail) {
			return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "EMAIL_ALREADY_EXISTS", Message: "Email já cadastrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: updated.ID, Name: updated.Name, Email: updated.Email, Balance: updated.Balance}
	return c.JSON(out)
}

// DeleteCustomer godoc
// @Summary      Deleta um usuário pelo ID
// @Description  Endpoint para deletar um usuário pelo ID
// @Tags         Clientes
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID do usuário"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [delete]
func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.Delete(c.UserContext(), id); err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// DepositCustomer godoc
// @Summary      Deposita um valor na conta do usuário
// @Description  Endpoint para depositar um valor na conta do usuário
// @Tags         Transações
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "ID do usuário"
// @Param        deposit  body      types.TransactionRequest  true  "Dados do depósito"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/depositar [post]
func (h *CustomerHandler) Deposit(c *fiber.Ctx) error {
	id := c.Params("id")

	req := &types.TransactionRequest{}
	err := req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	cust, err := h.service.Transactions(c.UserContext(), id, req.Amount)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INSUFICIENT_BALANCE", Message: "Saldo insuficiente"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Name: cust.Name, Email: cust.Email, Balance: cust.Balance}
	return c.JSON(out)
}

// WithdrawCustomer godoc
// @Summary      Saca um valor da conta do usuário
// @Description  Endpoint para sacar um valor da conta do usuário
// @Tags         Transações
// @Accept       json
// @Produce      json
// @Param        id        path      string  true  "ID do usuário"
// @Param        withdraw  body      types.TransactionRequest  true  "Dados do saque"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/sacar [post]
func (h *CustomerHandler) Withdraw(c *fiber.Ctx) error {
	id := c.Params("id")

	req := &types.TransactionRequest{}
	err := req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	cust, err := h.service.Transactions(c.UserContext(), id, req.Amount.Neg())
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INSUFICIENT_BALANCE", Message: "Saldo insuficiente"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Name: cust.Name, Email: cust.Email, Balance: cust.Balance}
	return c.JSON(out)
}

// GetTransactions retorna o histórico de transações de um cliente com paginação.
//
// GetTransactions godoc
// @Summary      Lista todas as transações de um usuário
// @Description  Endpoint para listar o histórico de transações de um cliente, com suporte a paginação via query params `page` e `size`.
// @Tags         Transações
// @Accept       json
// @Produce      json
// @Param        id    path      string true  "ID do usuário (UUID)"
// @Param        page  query     int    false "Número da página (default: 1)"
// @Param        size  query     int    false "Itens por página (default: 10)"
// @Success      200  {object}  map[string]interface{}  "Retorna metadados de paginação e a lista de transações"
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/transacoes [get]
func (h *CustomerHandler) GetTransactions(c *fiber.Ctx) error {
	id := c.Params("id")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	if _, err := h.service.GetByID(c.UserContext(), id); err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(
				types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"},
			)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(
			types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()},
		)
	}

	txs, total, err := h.service.ListTransactions(c.UserContext(), id, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()},
		)
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))

	if page > totalPages && totalPages != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{
			Code:    "INVALID_PAGE",
			Message: fmt.Sprintf("A página %d não existe. Total de páginas: %d", page, totalPages),
		})
	}

	out := make([]types.TransactionDto, 0, len(txs))
	for _, t := range txs {
		out = append(out, types.TransactionDto{
			TransactionID: t.TransactionID,
			CustomerID:    t.CustomerID,
			Amount:        t.Amount,
			Type:          t.Type,
			CreatedAt:     t.CreatedAt,
		})
	}

	return c.JSON(fiber.Map{
		"page":        page,
		"size":        size,
		"total_items": total,
		"total_pages": totalPages,
		"items":       out,
	})
}
