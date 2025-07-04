package https

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/internal-transfer/domain/banking"
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/pkg/ctx"
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
	api.POST("/transactions", handler.Transaction)
}

func (h *bankingHandler) CreateAccount(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)
	var account dto.AccountCreationRequest
	if err := ac.CustomBind(&account); err != nil {
		return ac.CustomResponse("Bad Request", nil, "", "Invalid request", http.StatusBadRequest, nil)
	}
	if err := h.usecase.CreateAccount(c, account.AccountID, account.InitialBalance); err != nil {
		return ac.CustomResponse("Internal Server Error", nil, "", "Failed to create account", http.StatusInternalServerError, nil)
	}
	return ac.CustomResponse("Success", nil, "Account created successfully", "", http.StatusCreated, nil)
}

func (h *bankingHandler) GetAccount(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)
	accountID := c.Param("id")
	if accountID == "" {
		return ac.CustomResponse("Bad Request", nil, "", "Account ID is required", http.StatusBadRequest, nil)
	}
	id, err := strconv.Atoi(accountID)
	if err != nil {
		return ac.CustomResponse("Bad Request", nil, "", "Invalid account ID format", http.StatusBadRequest, nil)
	}
	account, err := h.usecase.GetAccount(id)
	if err != nil {
		return ac.CustomResponse("Internal Server Error", nil, "", "Failed to retrieve account", http.StatusInternalServerError, nil)
	}
	return ac.CustomResponse("Success", account, "Account retrieved successfully", "", http.StatusOK, nil)
}

func (h *bankingHandler) Transaction(c echo.Context) error {
	ac := c.(*ctx.CustomApplicationContext)
	var transaction dto.TransactionRequest
	if err := ac.CustomBind(&transaction); err != nil {
		return ac.CustomResponse("Bad Request", nil, "", "Invalid request body", http.StatusBadRequest, nil)
	}
	if err := h.usecase.Transaction(c, transaction.SourceAccountID, transaction.DestinationAccountID,
		transaction.Amount); err != nil {
		return ac.CustomResponse("Internal Server Error", nil, "", "Transaction failed: "+err.Error(), http.StatusInternalServerError, nil)
	}
	return ac.CustomResponse("Success", nil, "Transaction completed successfully", "", http.StatusOK, nil)
}
