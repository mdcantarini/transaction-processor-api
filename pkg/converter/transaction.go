package converter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/mdcantarini/transaction-processor-api/pkg/model"
)

var CSVRecordsToTransactions = csvRecordsToTransactions

const (
	convertingError = `error converting`
)

func csvRecordsToTransactions(records [][]string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	for _, record := range records {
		transactionID, err := stringToInt(record[0])
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s %s", convertingError, "transactionID"))
		}

		date, err := stringToDate(record[1])
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s %s", convertingError, "date"))
		}

		transactionAmount, err := stringToDecimal(record[2])
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s %s", convertingError, "transactionAmount"))
		}

		accountID, err := stringToInt(record[3])
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("%s %s", convertingError, "accountID"))
		}

		transactions = append(transactions, model.Transaction{
			TransactionID:     transactionID,
			Date:              date,
			TransactionAmount: transactionAmount,
			AccountID:         accountID,
		})
	}

	return transactions, nil
}

func stringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("unable to convert %s into int", s)
	}

	return i, nil
}

func stringToDecimal(s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, fmt.Errorf("unable to convert %s into decimal", s)
	}

	return d, nil
}

const dateFormat = "2006-01-02"

func stringToDate(s string) (time.Time, error) {
	d, err := time.Parse(dateFormat, s)
	if err != nil {
		return time.Now(), fmt.Errorf("unable to convert %s into decimal", s)
	}

	return d, nil
}

func TransactionsToEmailTemplate(transactions []model.Transaction) string {
	debitTotal := decimal.Zero
	debitCount := decimal.Zero
	creditTotal := decimal.Zero
	creditCount := decimal.Zero
	transactionCountByMonth := map[string]int{}

	for _, t := range transactions {
		if t.TransactionAmount.IsPositive() {
			creditCount = creditCount.Add(decimal.RequireFromString("1"))
			creditTotal = creditTotal.Add(t.TransactionAmount)
		} else {
			debitCount = debitCount.Add(decimal.RequireFromString("1"))
			debitTotal = debitTotal.Add(t.TransactionAmount)
		}

		if _, ok := transactionCountByMonth[t.Date.Month().String()]; !ok {
			transactionCountByMonth[t.Date.Month().String()] = 0
		}

		transactionCountByMonth[t.Date.Month().String()] = transactionCountByMonth[t.Date.Month().String()] + 1
	}

	bodyRows := []string{
		"<html>",
		`<body style="font-family: Verdana, sans-serif; margin: 0; padding: 0;">`,
		"<h3>Historical summary</h3>",
		fmt.Sprintf("<p>Total balance is: %s</p>", creditTotal.Add(debitTotal).String()),
		fmt.Sprintf("<p>Average debit amount: %s</p>", debitTotal.Add(debitCount).String()),
		fmt.Sprintf("<p>Average credit amount: %s</p>", creditTotal.Div(creditCount).String()),
	}

	bodyRows = append(bodyRows, "<h4>Monthly summary</h4>")
	bodyRows = append(bodyRows, "<ul>")
	for month, count := range transactionCountByMonth {
		bodyRows = append(bodyRows, fmt.Sprintf("<li>Number of transactions in %s: %d</li>", month, count))
	}
	bodyRows = append(bodyRows, "</ul>")

	bodyRows = append(bodyRows, "<p>You will find the latest processed report attached to this email<p>")

	bodyRows = append(bodyRows, `<hr>`)
	bodyRows = append(bodyRows, fmt.Sprintf(`<p><img src="%s" alt="logo"><p>`, os.Getenv("EMAIL_LOGO_URL")))
	bodyRows = append(bodyRows, "</body>")
	bodyRows = append(bodyRows, "</html>")

	return strings.Join(bodyRows, "")
}
