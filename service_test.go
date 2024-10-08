package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mdcantarini/transaction-processor-api/pkg/model"
	"github.com/mdcantarini/transaction-processor-api/pkg/utils"
)

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) UpsertTransactions(transactions []model.Transaction) error {
	args := m.Called(transactions)
	return args.Error(0)
}

type MockAccountRepo struct {
	mock.Mock
}

func (m *MockAccountRepo) GetAccount(accountID int) (*model.Account, error) {
	args := m.Called(accountID)
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepo) UpsertAccounts(accounts []model.Account) error {
	args := m.Called(accounts)
	return args.Error(0)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to, subject, body string, attachments ...string) error {
	args := m.Called(to, subject, body, attachments)
	return args.Error(0)
}

func TestRunDailyReport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// mock success response in utils.ParseCSVFile function
	mockParseCSVFile(t, [][]string{{"1", "2023-12-15", "60.5", "1"}, {"2", "2023-12-15", "-10.3", "1"}}, nil)

	// mock transaction repo
	mockTransactionRepo := new(MockTransactionRepo)
	mockTransactionRepo.On("UpsertTransactions", mock.Anything).
		Return(nil).
		Times(1)

	// mock account repo
	mockAccountRepo := new(MockAccountRepo)
	mockAccountRepo.On("GetAccount", mock.Anything).
		Return(&model.Account{ID: 1, Email: "test@email.com"}, nil).
		Times(1)

	// mock email sender
	mockEmailSender := new(MockEmailSender)
	mockEmailSender.On("SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Times(1)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	service := &Service{
		accountRepo:     mockAccountRepo,
		transactionRepo: mockTransactionRepo,
		emailSender:     mockEmailSender,
	}
	service.RunDailyReport(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func mockParseCSVFile(t *testing.T, expectedRes [][]string, expectedErr error) {
	orig := utils.ParseCSVFile
	utils.ParseCSVFile = func(filepath string) ([][]string, error) {
		return expectedRes, expectedErr
	}
	t.Cleanup(func() {
		utils.ParseCSVFile = orig
	})
}
