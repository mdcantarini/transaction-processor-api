package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mdcantarini/transaction-processor-api/pkg/converter"
	"github.com/mdcantarini/transaction-processor-api/pkg/model"
	"github.com/mdcantarini/transaction-processor-api/pkg/utils"
	"github.com/mdcantarini/transaction-processor-api/pkg/utils/email"
	"github.com/mdcantarini/transaction-processor-api/pkg/utils/email/implementation"
)

type Service struct {
	transactionRepo model.ITransaction
	accountRepo     model.IAccount
	emailSender     email.EmailSender
}

func NewService() *Service {
	db := utils.MustCreateDBConnection()

	return &Service{
		transactionRepo: model.TransactionRepository{DB: db},
		accountRepo:     model.AccountRepository{DB: db},
		emailSender: implementation.Mailtrap{
			FromEmail: os.Getenv("MAILTRAP_FROM_EMAIL"),
			Host:      os.Getenv("MAILTRAP_HOST"),
			Token:     os.Getenv("MAILTRAP_TOKEN"),
		},
	}
}

func (s *Service) RunDailyReport(c *gin.Context) {
	filePath := os.Getenv("TRANSACTIONS_FILE_PATH")

	records, err := utils.ParseCSVFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to parse csv file: %w", err))
	}

	transactions, err := converter.CSVRecordsToTransactions(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to convert csv records to transactions: %w", err))
	}

	err = s.transactionRepo.UpsertTransactions(transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to insert transactions: %w", err))
	}

	err = s.sendDailyReportByEmail(transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to send daily report by email: %w", err))
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Service) sendDailyReportByEmail(txs []model.Transaction) error {
	transactionsByAccount := groupTransactionsByAccount(txs)

	for accountID, transactions := range transactionsByAccount {
		account, err := s.accountRepo.GetAccount(accountID)
		if err != nil {
			return fmt.Errorf("unable to fetch user %d", accountID)
		}

		to := account.Email
		subject := fmt.Sprintf("Daily report for Account %d", accountID)
		body := converter.TransactionsToEmailTemplate(transactions)
		attachment := os.Getenv("TRANSACTIONS_FILE_PATH")

		err = s.emailSender.SendEmail(to, subject, body, attachment)
		if err != nil {
			fmt.Println("unable to send email:", err.Error())
			return err
		}
	}

	return nil
}

func groupTransactionsByAccount(txs []model.Transaction) map[int][]model.Transaction {
	result := map[int][]model.Transaction{}

	for _, tx := range txs {
		if _, ok := result[tx.AccountID]; !ok {
			result[tx.AccountID] = []model.Transaction{}
		}

		result[tx.AccountID] = append(result[tx.AccountID], tx)
	}

	return result
}
