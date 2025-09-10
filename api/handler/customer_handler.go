package handler

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"

	"case-itau/api/types"
	repository "case-itau/repository/interface"
	"case-itau/services/customer"
)

// CustomerHandler holds dependencies
type CustomerHandler struct {
	service  *customer.Service
	validate *validator.Validate
}

// NewCustomerHandler creates a handler
func NewCustomerHandler(s *customer.Service) *CustomerHandler {
	return &CustomerHandler{
		service:  s,
		validate: validator.New(),
	}
}

// GetCustomers godoc
// @Summary      Lista todos os usuários
// @Description  Endpoint para listar todos os usuários
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes [get]
func (h *CustomerHandler) List(c *fiber.Ctx) error {
	list, err := h.service.ListAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: "internal error"})
	}

	out := make([]types.CustomerDto, 0, len(list))
	for _, it := range list {
		out = append(out, types.CustomerDto{
			ID: it.ID, Nome: it.Nome, Email: it.Email, Saldo: it.Saldo,
		})
	}
	return c.JSON(out)
}

// GetCustomer godoc
// @Summary      Obtém um usuário pelo ID
// @Description  Endpoint para obter um usuário pelo ID
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do usuário"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [get]
func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid id"})
	}
	cust, err := h.service.GetByID(uint(id64))
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Message: "cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: "internal error"})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
	return c.JSON(out)
}

// CreateCustomer godoc
// @Summary      Cria um novo usuário
// @Description  Endpoint para criar um novo usuário
// @Accept       json
// @Produce      json
// @Param        customer  body      types.CreateCustomerRequest  true  "Dados do usuário"
// @Success      201  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes [post]
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req types.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid json"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "validation error"})
	}

	model := repository.Clientes{
		Nome:  req.Nome,
		Email: req.Email,
		Saldo: decimal.Zero,
	}
	created, err := h.service.Create(model)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: "could not create"})
	}
	out := types.CustomerDto{ID: created.ID, Nome: created.Nome, Email: created.Email, Saldo: created.Saldo}

	return c.Status(fiber.StatusCreated).JSON(out)
}

// UpdateCustomer godoc
// @Summary      Atualiza um usuário existente
// @Description  Endpoint para atualizar um usuário existente
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "ID do usuário"
// @Param        customer  body      types.UpdateCustomerRequest  true  "Dados do usuário"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [put]
func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid id"})
	}

	var req types.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid json"})
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "validation error"})
	}

	in := repository.Clientes{Nome: req.Nome, Email: req.Email}
	updated, err := h.service.Update(uint(id64), in)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Message: "cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: "could not update"})
	}
	out := types.CustomerDto{ID: updated.ID, Nome: updated.Nome, Email: updated.Email, Saldo: updated.Saldo}
	return c.JSON(out)
}

// DeleteCustomer godoc
// @Summary      Deleta um usuário pelo ID
// @Description  Endpoint para deletar um usuário pelo ID
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID do usuário"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [delete]
func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid id"})
	}
	if err := h.service.Delete(uint(id64)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: "could not delete"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// DepositCustomer godoc
// @Summary      Deposita um valor na conta do usuário
// @Description  Endpoint para depositar um valor na conta do usuário
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "ID do usuário"
// @Param        deposit  body      types.TransactionRequest  true  "Dados do depósito"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/depositar [post]
func (h *CustomerHandler) Deposit(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid id"})
	}

	var req types.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid json"})
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "validation error"})
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "amount must be > 0"})
	}

	cust, err := h.service.Deposit(uint(id64), req.Amount)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Message: "cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
	return c.JSON(out)
}

// WithdrawCustomer godoc
// @Summary      Saca um valor da conta do usuário
// @Description  Endpoint para sacar um valor da conta do usuário
// @Accept       json
// @Produce      json
// @Param        id        path      int  true  "ID do usuário"
// @Param        withdraw  body      types.TransactionRequest  true  "Dados do saque"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/sacar [post]
func (h *CustomerHandler) Withdraw(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid id"})
	}

	var req types.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "invalid json"})
	}
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Message: "amount must be > 0"})
	}

	cust, err := h.service.Withdraw(uint(id64), req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
	return c.JSON(out)
}
