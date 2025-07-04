package banking

import (
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/models"
)

type Usecase interface {
	CreateAccount(accountID int, balance string) error
	GetAccount(accountID int) (dto.AccountResponse, error)
	Transection(fromAccountID int, toAccountID int, amount string) error
}
type Repository interface {
	CreateAccount(account models.Account) error
	GetAccount(accountID int) (models.Account, error)
	UpdateAccount(account models.Account) error
	Transection(transaction models.Transaction) error
}
