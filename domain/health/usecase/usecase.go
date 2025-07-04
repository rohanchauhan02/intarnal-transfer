package usecase

import "github.com/rohanchauhan02/internal-transfer/domain/health"

type healthUsecase struct {
	repo health.Repository
}

// NewHealthUsecase creates a new Health usecase instance
func NewHealthUsecase(repo health.Repository) health.Usecase {
	return &healthUsecase{
		repo: repo,
	}
}
