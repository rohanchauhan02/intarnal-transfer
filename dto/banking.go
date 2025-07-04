package dto


type AccountCreationRequest struct {
	AccountID      int    `json:"account_id"`
	InitialBalance string `json:"initial_balance"`
}

type AccountResponse struct {
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

type TransactionRequest struct {
	SourceAccountID      int    `json:"source_account_id"`
	DestinationAccountID int    `json:"destination_account_id"`
	Amount               string `json:"amount"`
}
