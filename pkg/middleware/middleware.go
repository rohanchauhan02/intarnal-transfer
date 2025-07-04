package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func MiddlewareRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Request().Header.Get(echo.HeaderXRequestID)
			if requestID == "" {
				requestID = generateRequestID()
			}
			c.Request().Header.Set(echo.HeaderXRequestID, requestID)
			c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			return next(c)
		}

	}
}

func generateRequestID() string {
	return uuid.New().String()
}
