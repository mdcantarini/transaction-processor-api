package converter

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestCSVRecordsToTransactions_Success(t *testing.T) {
	records := [][]string{{"1", "2023-12-15", "60.5", "1"}, {"2", "2023-12-16", "-10.3", "2"}}

	transactions, err := CSVRecordsToTransactions(records)
	require.NoError(t, err)

	require.Equal(t, transactions[0].TransactionID, 1)
	date, err := time.Parse(dateFormat, "2023-12-15")
	require.NoError(t, err)
	require.Equal(t, transactions[0].Date, date)
	require.Equal(t, transactions[0].TransactionAmount, decimal.RequireFromString("60.5"))
	require.Equal(t, transactions[0].AccountID, 1)

	require.Equal(t, transactions[1].TransactionID, 2)
	date, err = time.Parse(dateFormat, "2023-12-16")
	require.NoError(t, err)
	require.Equal(t, transactions[1].Date, date)
	require.Equal(t, transactions[1].TransactionAmount, decimal.RequireFromString("-10.3"))
	require.Equal(t, transactions[1].AccountID, 2)
}

func TestCSVRecordsToTransactions_ErrorConvertingTransactionID(t *testing.T) {
	records := [][]string{{"a", "2023-12-15", "60.5", "1"}}

	_, err := CSVRecordsToTransactions(records)
	require.ErrorContains(t, err, fmt.Sprintf("%s %s", convertingError, "transactionID"))
}

func TestCSVRecordsToTransactions_ErrorConvertingTransactionDate(t *testing.T) {
	records := [][]string{{"1", "202a-12-15", "60.5", "1"}}

	_, err := CSVRecordsToTransactions(records)
	require.ErrorContains(t, err, fmt.Sprintf("%s %s", convertingError, "date"))
}

func TestCSVRecordsToTransactions_ErrorConvertingTransactionAmount(t *testing.T) {
	records := [][]string{{"1", "2023-12-15", "60.a", "1"}}

	_, err := CSVRecordsToTransactions(records)
	require.ErrorContains(t, err, fmt.Sprintf("%s %s", convertingError, "transactionAmount"))
}

func TestCSVRecordsToTransactions_ErrorConvertingAccountID(t *testing.T) {
	records := [][]string{{"1", "2023-12-15", "60.5", "a"}}

	_, err := CSVRecordsToTransactions(records)
	require.ErrorContains(t, err, fmt.Sprintf("%s %s", convertingError, "accountID"))
}
