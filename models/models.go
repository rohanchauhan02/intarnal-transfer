// This package contains the data models for the application.
package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	AccountID int    `json:"account_id"`
	Balance   string `json:"balance"`
}

type Transaction struct {
	gorm.Model
	SourceAccountID      int     `json:"source_account_id"`
	DestinationAccountID int     `json:"destination_account_id"`
	Amount               string  `json:"amount"`
}
