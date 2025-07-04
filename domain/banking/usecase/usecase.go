package usecase

import "github.com/rohanchauhan02/internal-transfer/domain/banking"

type bankingUsecase struct {
	repo banking.Repository
}

// NewBankingUsecase creates a new banking usecase instance
func NewBankingUsecase(repo banking.Repository) banking.Usecase {
	return &bankingUsecase{
		repo: repo,
	}
}
