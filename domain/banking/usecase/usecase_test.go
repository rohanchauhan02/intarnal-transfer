package usecase

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/rohanchauhan02/internal-transfer/dto"
	mock_banking "github.com/rohanchauhan02/internal-transfer/file/mocks/mock_banking"
	"github.com/rohanchauhan02/internal-transfer/models"
	"github.com/rohanchauhan02/internal-transfer/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBankingUsecase_CreateAccount(t *testing.T) {
	db, sqlmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		accountID     int
		balance       string
		mockSetup     func(repo *mock_banking.MockRepository)
		expectedError error
	}{
		{
			name:      "Create Account Success",
			accountID: 1,
			balance:   "1000.00",
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "Create Account Failure",
			accountID: 2,
			balance:   "500.00",
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(errors.New("failed to create account"))
			},
			expectedError: errors.New("failed to create account"),
		},
	}

	for _, tt := range tests {
		sqlmock.ExpectBegin()
		sqlmock.ExpectCommit()

		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_banking.NewMockRepository(ctrl)
			tt.mockSetup(mockRepo)

			usecase := NewBankingUsecase(mockRepo)

			c := &ctx.CustomApplicationContext{
				PostgresDB: gormDB,
			}

			err := usecase.CreateAccount(c, tt.accountID, tt.balance)
			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
func TestBankingUsecase_GetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_banking.NewMockRepository(ctrl)
	usecase := NewBankingUsecase(mockRepo)

	tests := []struct {
		name          string
		accountID     int
		mockSetup     func()
		expectedResp  dto.AccountResponse
		expectedError error
	}{
		{
			name:      "Get Account Success",
			accountID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAccount(1).
					Return(models.Account{AccountID: 1, Balance: "1000.00"}, nil)
			},
			expectedResp: dto.AccountResponse{
				AccountID: 1,
				Balance:   "1000.00",
			},
			expectedError: nil,
		},
		{
			name:      "Get Account Not Found",
			accountID: 2,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAccount(2).
					Return(models.Account{}, errors.New("account not found"))
			},
			expectedResp:  dto.AccountResponse{},
			expectedError: errors.New("account not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			resp, err := usecase.GetAccount(tt.accountID)
			if tt.expectedError != nil {
				assert.ErrorContains(t, err, tt.expectedError.Error())
				assert.Equal(t, tt.expectedResp, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
func TestBankingUsecase_Transaction(t *testing.T) {
	db, sqlmock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm database: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		fromAccountID int
		toAccountID   int
		amount        string
	}
	tests := []struct {
		name          string
		args          args
		mockSetup     func(repo *mock_banking.MockRepository)
		sqlSetup      func()
		expectedError string
	}{
		{
			name: "Transaction Success",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(nil).Times(2)
				repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).Return(nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectCommit()
			},
			expectedError: "",
		},
		{
			name: "Insufficient Balance",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "600.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "insufficient balance",
		},
		{
			name: "Invalid Sender Balance",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "invalid"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "invalid balance in sender's account",
		},
		{
			name: "Invalid Receiver Balance",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "invalid"}, nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "invalid balance in receiver's account",
		},
		{
			name: "Invalid Transfer Amount",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "invalid"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "invalid transfer amount",
		},
		{
			name: "GetAccountTx Sender Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{}, errors.New("sender not found"))
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "sender not found",
		},
		{
			name: "GetAccountTx Receiver Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{}, errors.New("receiver not found"))
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "receiver not found",
		},
		{
			name: "UpdateAccount Sender Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(errors.New("update sender error"))
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "update sender error",
		},
		{
			name: "UpdateAccount Receiver Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(errors.New("update receiver error"))
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "update receiver error",
		},
		{
			name: "Transaction Insert Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(nil).Times(2)
				repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).Return(errors.New("insert transaction error"))
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectRollback()
			},
			expectedError: "insert transaction error",
		},
		{
			name: "Commit Error",
			args: args{fromAccountID: 1, toAccountID: 2, amount: "100.00"},
			mockSetup: func(repo *mock_banking.MockRepository) {
				repo.EXPECT().GetAccountTx(gomock.Any(), 1).Return(models.Account{AccountID: 1, Balance: "500.00"}, nil)
				repo.EXPECT().GetAccountTx(gomock.Any(), 2).Return(models.Account{AccountID: 2, Balance: "200.00"}, nil)
				repo.EXPECT().UpdateAccount(gomock.Any(), gomock.AssignableToTypeOf(models.Account{})).Return(nil).Times(2)
				repo.EXPECT().Transaction(gomock.Any(), gomock.Any()).Return(nil)
			},
			sqlSetup: func() {
				sqlmock.ExpectBegin()
				sqlmock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedError: "failed to commit transaction: commit error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock_banking.NewMockRepository(ctrl)
			tt.mockSetup(mockRepo)
			tt.sqlSetup()

			usecase := NewBankingUsecase(mockRepo)
			c := &ctx.CustomApplicationContext{
				PostgresDB: gormDB,
			}

			err := usecase.Transaction(c, tt.args.fromAccountID, tt.args.toAccountID, tt.args.amount)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
