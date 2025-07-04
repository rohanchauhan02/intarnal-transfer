package usecase

import (
	"errors"

	"github.com/rohanchauhan02/internal-transfer/domain/health"
)

type healthUsecase struct {
	repo health.Repository
}

// NewHealthUsecase creates a new Health usecase instance
func NewHealthUsecase(repo health.Repository) health.Usecase {
	return &healthUsecase{
		repo: repo,
	}
}

// CheckHealth checks the health of the service
func (u *healthUsecase) CheckHealth() (map[string]string, error) {
	dbStatus, err := u.repo.PingDatabase()
	if err != nil {
		return map[string]string{
			"database": "unhealthy",
		}, errors.New("database is unreachable")
	}

	// If everything is fine, return a healthy response
	return map[string]string{
		"database": dbStatus,
	}, nil
}
