package https

import (
	"github.com/rohanchauhan02/internal-transfer/domain/health"
)

type healthHandler struct {
	usecase health.Usecase
}

// NewHealthHandler creates a new health handler with the provided usecase.
func NewHealthHandler(usecase health.Usecase) *healthHandler {
	return &healthHandler{
		usecase: usecase,
	}
}
