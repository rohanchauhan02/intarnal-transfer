package https

import "github.com/rohanchauhan02/internal-transfer/domain/banking"

type bankingHandler struct {
	usecase banking.Usecase
}

// NewBankingHandler creates a new banking handler with the provided usecase.
func NewBankingHandler(usecase banking.Usecase) *bankingHandler {
	return &bankingHandler{
		usecase: usecase,
	}
}
