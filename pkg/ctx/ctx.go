package ctx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/utils"

	"gorm.io/gorm"
)

// Validator wraps the Go Playground validator
type Validator struct {
	Validator *validator.Validate
}

// Validate performs validation on struct input
func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

// CustomApplicationContext extends Echo's Context to include additional dependencies.
type CustomApplicationContext struct {
	echo.Context
	PostgresDB *gorm.DB
}

// CustomResponse formats and sends a structured JSON response.
func (c *CustomApplicationContext) CustomResponse(status string, data any, message string, errMsg string, code int, meta any) error {
	response := &dto.ResponsePattern{
		RequestID:    c.Request().Header.Get(echo.HeaderXRequestID),
		Status:       status,
		Data:         data,
		Message:      message,
		ErrorMessage: errMsg,
		Code:         code,
		Meta:         meta,
	}

	respBytes, err := json.Marshal(response)
	if err == nil {
		log.Infof("%s -- RESPONSE -- %s", utils.GetCallerMethod(), string(respBytes))
	} else {
		log.Errorf("Failed to marshal response: %v", err)
	}

	return c.JSON(code, response)
}

// CustomBind binds and validates incoming request data.
func (c *CustomApplicationContext) CustomBind(i any) error {
	if err := c.Bind(i); err != nil {
		log.Warnf("%s -- Failed to bind request payload: %v", utils.GetCallerMethod(), err)
		return err
	}

	if err := c.Validate(i); err != nil {
		log.Warnf("%s -- Validation failed: %v", utils.GetCallerMethod(), err)
		return mapValidationErrors(err)
	}

	reqBytes, err := json.Marshal(i)
	if err == nil {
		log.Infof("%s -- Payload -- %s", utils.GetCallerMethod(), string(reqBytes))
	} else {
		log.Errorf("Failed to marshal request payload: %v", err)
	}

	return nil
}

// mapValidationErrors converts validation errors into a user-friendly format.
func mapValidationErrors(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var errorMessages []string
		for _, e := range validationErrs {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", e.Field()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}
		return errors.New(strings.Join(errorMessages, "; "))
	}
	return err
}
