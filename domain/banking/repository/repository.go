package repository

import (
	"github.com/rohanchauhan02/internal-transfer/domain/banking"
	"github.com/rohanchauhan02/internal-transfer/models"
	"gorm.io/gorm"
)

type bankingRepository struct {
	db *gorm.DB
}

// NewBankingRepository creates a new Repository instance
func NewBankingRepository(db *gorm.DB) banking.Repository {
	return &bankingRepository{
		db: db,
	}
}

// CreateAccount creates a new account in the database
func (r *bankingRepository) CreateAccount(tx *gorm.DB, account models.Account) error {
	return tx.Create(&account).Error
}

// GetAccount retrieves an account by its ID
func (r *bankingRepository) GetAccount(accountID int) (models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_id", accountID).Find(&account).Error; err != nil {
		return models.Account{}, err
	}
	return account, nil
}

// GetAccountTx retrieves an account by its ID within a transaction
func (r *bankingRepository) GetAccountTx(tx *gorm.DB, accountID int) (models.Account, error) {
	var account models.Account
	if err := tx.Raw("SELECT * FROM accounts WHERE account_id = ? FOR UPDATE", accountID).Scan(&account).Error; err != nil {
		return models.Account{}, err
	}
	if err := tx.Where("account_id = ?", accountID).Find(&account).Error; err != nil {
		return models.Account{}, err
	}
	return account, nil
}

// UpdateAccount updates an existing account in the database
func (r *bankingRepository) UpdateAccount(tx *gorm.DB, account models.Account) error {
	return tx.Save(&account).Error
}

// Transaction processes a transaction between accounts
func (r *bankingRepository) Transaction(tx *gorm.DB, transaction models.Transaction) error {
	return tx.Create(&transaction).Error
}
