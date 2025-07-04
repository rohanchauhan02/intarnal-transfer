package dto

type AccountCreationRequest struct {
	AccountID      int    `json:"account_id" validate:"required"`
	InitialBalance string `json:"initial_balance" validate:"required"`
}

type AccountResponse struct {
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

type TransactionRequest struct {
	SourceAccountID      int    `json:"source_account_id" validate:"required"`
	DestinationAccountID int    `json:"destination_account_id" validate:"required"`
	Amount               string `json:"amount" validate:"required"`
}
