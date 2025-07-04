package repository

import (
	"github.com/rohanchauhan02/internal-transfer/domain/banking"
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
