package repository

import (
	"context"
	"errors"

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

// CheckHealth checks the health of the service
func (r *healthRepository) PingDatabase() (string, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		return "unhealthy", errors.New("failed to get database instance")
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		return "unhealthy", errors.New("database is unreachable")
	}

	return "healthy", nil
}
