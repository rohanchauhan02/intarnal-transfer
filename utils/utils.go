package utils

import (
	"path"
	"runtime"
	"strings"

	"github.com/go-playground/validator"
)

// CustomValidator return custom validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate will validate given input with related struct
func (cv *CustomValidator) Validate(i any) error {
	return cv.Validator.Struct(i)
}

// func DefaultValidator function to give difault validation all incoming request
func DefaultValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func GetCallerMethod() string {
	var source string
	if pc, _, _, ok := runtime.Caller(2); ok {
		var funcName string
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = fn.Name()
			if i := strings.LastIndex(funcName, "."); i != -1 {
				funcName = funcName[i+1:]
			}
		}
		source = path.Base(funcName)
	}
	return source
}
