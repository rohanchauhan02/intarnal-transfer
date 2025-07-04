package usecase

import (
	"errors"
	"strconv"

	"github.com/rohanchauhan02/internal-transfer/domain/banking"
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/models"
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
func (u *bankingUsecase) CreateAccount(accountID int, balance string) error {
	account := models.Account{
		AccountID: accountID,
		Balance:   balance,
	}
	return u.repo.CreateAccount(account)
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

// Transection transfers funds between accounts
func (u *bankingUsecase) Transection(fromAccountID int, toAccountID int, amount string) error {
	fromAccount, err := u.repo.GetAccount(fromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := u.repo.GetAccount(toAccountID)
	if err != nil {
		return err
	}

	// Convert balances and amount to float64 for arithmetic
	fromBalance, err := strconv.ParseFloat(fromAccount.Balance, 64)
	if err != nil {
		return errors.New("invalid balance in sender's account")
	}
	toBalance, err := strconv.ParseFloat(toAccount.Balance, 64)
	if err != nil {
		return errors.New("invalid balance in receiver's account")
	}
	transferAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return errors.New("invalid transfer amount")
	}

	// Check if sufficient balance is available
	if fromBalance < transferAmount {
		return errors.New("insufficient balance")
	}

	// Deduct amount from sender's account
	fromBalance -= transferAmount
	fromAccount.Balance = strconv.FormatFloat(fromBalance, 'f', -1, 64)
	if err := u.repo.UpdateAccount(fromAccount); err != nil {
		return err
	}

	// Add amount to receiver's account
	toBalance += transferAmount
	toAccount.Balance = strconv.FormatFloat(toBalance, 'f', -1, 64)
	return u.repo.UpdateAccount(toAccount)
}
