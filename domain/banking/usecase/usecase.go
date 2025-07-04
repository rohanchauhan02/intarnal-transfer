package usecase

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/internal-transfer/domain/banking"
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/models"
	"github.com/rohanchauhan02/internal-transfer/pkg/ctx"
	"github.com/shopspring/decimal"
)

type bankingUsecase struct {
	repo banking.Repository
}

// NewBankingUsecase creates a new banking usecase instance
func NewBankingUsecase(repo banking.Repository) banking.Usecase {
	return &bankingUsecase{
		repo: repo,
	}
}

// CreateAccount creates a new account
func (u *bankingUsecase) CreateAccount(c echo.Context, accountID int, balance string) error {
	ac := c.(*ctx.CustomApplicationContext)

	// Check if account already exists
	existingAccount, err := u.repo.GetAccount(accountID)
	if err == nil && existingAccount.AccountID == accountID {
		return errors.New("account already exists with this user ID")
	}

	account := models.Account{
		AccountID: accountID,
		Balance:   balance,
	}
	tx := ac.PostgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errors.New("failed to start transaction")
	}
	if err := u.repo.CreateAccount(tx, account); err != nil {
		tx.Rollback()
		return errors.New("failed to create account: " + err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}
	return nil
}

// GetAccount retrieves account details by account ID
func (u *bankingUsecase) GetAccount(accountID int) (dto.AccountResponse, error) {
	account, err := u.repo.GetAccount(accountID)
	if err != nil {
		return dto.AccountResponse{}, err
	}
	return dto.AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance,
	}, nil
}

// Transaction transfers funds between accounts
func (u *bankingUsecase) Transaction(c echo.Context, fromAccountID int, toAccountID int, amount string) error {
	ac := c.(*ctx.CustomApplicationContext)
	tx := ac.PostgresDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errors.New("failed to start transaction")
	}

	fromAccount, err := u.repo.GetAccountTx(tx, fromAccountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	toAccount, err := u.repo.GetAccountTx(tx, toAccountID)
	if err != nil {
		tx.Rollback()
		return err
	}

	fromBalance, err := decimal.NewFromString(fromAccount.Balance)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid balance in sender's account")
	}
	toBalance, err := decimal.NewFromString(toAccount.Balance)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid balance in receiver's account")
	}
	transferAmount, err := decimal.NewFromString(amount)
	if err != nil {
		tx.Rollback()
		return errors.New("invalid transfer amount")
	}

	if fromBalance.LessThan(transferAmount) {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	fromAccount.Balance = fromBalance.Sub(transferAmount).String()
	if err := u.repo.UpdateAccount(tx, fromAccount); err != nil {
		tx.Rollback()
		return err
	}

	toAccount.Balance = toBalance.Add(transferAmount).String()
	if err := u.repo.UpdateAccount(tx, toAccount); err != nil {
		tx.Rollback()
		return err
	}

	transaction := models.Transaction{
		SourceAccountID:      fromAccountID,
		DestinationAccountID: toAccountID,
		Amount:               amount,
	}

	if err := u.repo.Transaction(tx, transaction); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}
