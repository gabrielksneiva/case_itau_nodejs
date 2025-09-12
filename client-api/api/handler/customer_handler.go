package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
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
			ID: it.ID, Nome: it.Nome, Email: it.Email, Saldo: it.Saldo,
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
// @Param        id   path      int  true  "ID do usuário"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id} [get]
func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_CUSTOMER_ID", Message: "Id de cliente inválido"})
	}
	cust, err := h.service.GetByID(c.UserContext(), id64)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
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

	model := repo.Clientes{
		Nome:  req.Nome,
		Email: req.Email,
		Saldo: decimal.Zero,
	}
	created, err := h.service.Create(c.UserContext(), model)
	if err != nil {
		if errors.Is(err, customer.ErrUniqueEmail) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "EMAIL_ALREADY_EXISTS", Message: "Email já registrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: created.ID, Nome: created.Nome, Email: created.Email, Saldo: created.Saldo}

	return c.Status(fiber.StatusCreated).JSON(out)
}

// UpdateCustomer godoc
// @Summary      Atualiza um usuário existente
// @Description  Endpoint para atualizar um usuário existente
// @Tags         Clientes
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
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_COSTUMER_ID", Message: "Id de cliente inválido"})
	}

	req := &types.CreateCustomerRequest{}
	err = req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	in := repo.Clientes{Nome: req.Nome, Email: req.Email}
	updated, err := h.service.Update(c.UserContext(), id64, in)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrUniqueEmail) {
			return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "EMAIL_ALREADY_EXISTS", Message: "Email já cadastrado"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: updated.ID, Nome: updated.Nome, Email: updated.Email, Saldo: updated.Saldo}
	return c.JSON(out)
}

// DeleteCustomer godoc
// @Summary      Deleta um usuário pelo ID
// @Description  Endpoint para deletar um usuário pelo ID
// @Tags         Clientes
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
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_COSTUMER_ID", Message: "Id de cliente inválido"})
	}
	if err := h.service.Delete(c.UserContext(), id64); err != nil {
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
// @Param        id        path      int  true  "ID do usuário"
// @Param        deposit  body      types.TransactionRequest  true  "Dados do depósito"
// @Success      200  {object}  types.CustomerDto
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /clientes/{id}/depositar [post]
func (h *CustomerHandler) Deposit(c *fiber.Ctx) error {
	id64, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_COSTUMER_ID", Message: "Id de cliente inválido"})
	}

	req := &types.TransactionRequest{}
	err = req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	cust, err := h.service.Transactions(c.UserContext(), id64, req.Valor)
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INSUFICIENT_BALANCE", Message: "Saldo insuficiente"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
	return c.JSON(out)
}

// WithdrawCustomer godoc
// @Summary      Saca um valor da conta do usuário
// @Description  Endpoint para sacar um valor da conta do usuário
// @Tags         Transações
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
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_CUSTOMER_ID", Message: "Id de cliente inválido"})
	}

	req := &types.TransactionRequest{}
	err = req.FromBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&types.ErrorResponse{Code: "INVALID_REQUEST", Message: "Json inválido"})
	}
	if err := req.IsValid(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INVALID_REQUEST", Message: err.Error()})
	}

	cust, err := h.service.Transactions(c.UserContext(), id64, req.Valor.Neg())
	if err != nil {
		if errors.Is(err, customer.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(types.ErrorResponse{Code: "CUSTOMER_NOT_FOUND", Message: "Cliente não encontrado"})
		}
		if errors.Is(err, customer.ErrInsufficientFunds) {
			return c.Status(fiber.StatusBadRequest).JSON(types.ErrorResponse{Code: "INSUFICIENT_BALANCE", Message: "Saldo insuficiente"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(types.ErrorResponse{Code: "INTERNAL_ERROR", Message: err.Error()})
	}
	out := types.CustomerDto{ID: cust.ID, Nome: cust.Nome, Email: cust.Email, Saldo: cust.Saldo}
	return c.JSON(out)
}
