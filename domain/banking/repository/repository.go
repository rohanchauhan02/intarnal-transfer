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
func (r *bankingRepository) CreateAccount(account models.Account) error {
	return r.db.Create(&account).Error
}

// GetAccount retrieves an account by its ID
func (r *bankingRepository) GetAccount(accountID int) (models.Account, error) {
	var account models.Account
	if err := r.db.First(&account, accountID).Error; err != nil {
		return models.Account{}, err
	}
	return account, nil
}

// UpdateAccount updates an existing account in the database
func (r *bankingRepository) UpdateAccount(account models.Account) error {
	return r.db.Save(&account).Error
}
