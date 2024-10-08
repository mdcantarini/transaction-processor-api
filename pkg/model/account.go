package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Account struct {
	ID    int
	Email string `gorm:"unique"`
}

type IAccount interface {
	GetAccount(accountID int) (*Account, error)
	UpsertAccounts([]Account) error
}

type AccountRepository struct {
	DB *gorm.DB
}

func (ar AccountRepository) GetAccount(accountID int) (*Account, error) {
	account := Account{}
	err := ar.DB.First(&account, accountID).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (ar AccountRepository) UpsertAccounts(accounts []Account) error {
	return ar.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&accounts).Error
}
