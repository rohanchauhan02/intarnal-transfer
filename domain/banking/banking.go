package banking

import (
	"github.com/labstack/echo/v4"
	"github.com/rohanchauhan02/internal-transfer/dto"
	"github.com/rohanchauhan02/internal-transfer/models"
	"gorm.io/gorm"
)

type Usecase interface {
	CreateAccount(echo.Context, int, string) error
	GetAccount(int) (dto.AccountResponse, error)
	Transaction(echo.Context, int, int, string) error
}
type Repository interface {
	CreateAccount(*gorm.DB, models.Account) error
	GetAccount(int) (models.Account, error)
	GetAccountTx(*gorm.DB, int) (models.Account, error)
	UpdateAccount(*gorm.DB, models.Account) error
	Transaction(*gorm.DB, models.Transaction) error
}
