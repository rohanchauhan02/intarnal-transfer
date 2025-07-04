package https

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/internal-transfer/domain/banking"
	"github.com/rohanchauhan02/internal-transfer/dto"
)

type bankingHandler struct {
	usecase banking.Usecase
}

// NewBankingHandler creates a new banking handler with the provided usecase.
func NewBankingHandler(e *echo.Echo, usecase banking.Usecase) {
	handler := &bankingHandler{
		usecase: usecase,
	}

	api := e.Group("/api/v1")
	api.POST("/accounts", handler.CreateAccount)
	api.GET("/accounts/:id", handler.GetAccount)
	api.POST("/transactions", handler.Transection)
}

func (h *bankingHandler) CreateAccount(c echo.Context) error {
	// Implementation for creating an account
	var account dto.AccountCreationRequest
	if err := c.Bind(&account); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	if err := h.usecase.CreateAccount(account.AccountID, account.InitialBalance); err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to create account"})
	}
	return c.JSON(201, map[string]string{"message": "Account created successfully"})
}

func (h *bankingHandler) GetAccount(c echo.Context) error {
	// Implementation for getting account details
	accountID := c.Param("id")
	if accountID == "" {
		return c.JSON(400, map[string]string{"error": "Account ID is required"})
	}
	id, err := strconv.Atoi(accountID)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid Account ID"})
	}
	account, err := h.usecase.GetAccount(id)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "Account not found"})
	}
	return c.JSON(201, account)
}

func (h *bankingHandler) Transection(c echo.Context) error {
	// Implementation for transferring funds between accounts
	var transaction dto.TransactionRequest
	if err := c.Bind(&transaction); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	if err := h.usecase.Transection(transaction.SourceAccountID, transaction.DestinationAccountID,
		transaction.Amount); err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to process transaction"})
	}
	return c.JSON(200, map[string]string{"message": "Transaction successful"})
}
