package https

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/internal-transfer/domain/health"
)

type healthHandler struct {
	usecase health.Usecase
}

// NewHealthHandler creates a new health handler with the provided usecase.
func NewHealthHandler(e *echo.Echo, usecase health.Usecase) {
	handler := &healthHandler{
		usecase: usecase,
	}
	api := e.Group("/api/v1")
	api.GET("/healthz", handler.CheckHealth)
}

func (h *healthHandler) CheckHealth(c echo.Context) error {

	// Call the health check use case
	status, err := h.usecase.CheckHealth()
	if err != nil {
		log.Printf("Health check failed: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"data":    nil,
			"message": "Service is unhealthy",
			"error":   err.Error(),
		})
	}

	log.Println("Health check passed")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"data":    status,
		"message": "Service is healthy",
		"error":   "",
	})
}
