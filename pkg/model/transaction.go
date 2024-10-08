package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Transaction struct {
	gorm.Model
	TransactionID     int `gorm:"unique"`
	Date              time.Time
	TransactionAmount decimal.Decimal
	AccountID         int
	Account           Account
}

type ITransaction interface {
	UpsertTransactions(transactions []Transaction) error
}

type TransactionRepository struct {
	DB *gorm.DB
}

func (tr TransactionRepository) UpsertTransactions(transactions []Transaction) error {
	return tr.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&transactions).Error
}
