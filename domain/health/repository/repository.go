package repository

import (
	"github.com/rohanchauhan02/internal-transfer/domain/health"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

// NewHealthRepository creates a new Repository instance
func NewHealthRepository(db *gorm.DB) health.Repository {
	return &healthRepository{
		db: db,
	}
}
