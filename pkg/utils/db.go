package utils

import (
	"os"

	"github.com/cenkalti/backoff/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mdcantarini/transaction-processor-api/pkg/model"
)

func MustCreateDBConnection() *gorm.DB {
	// retry the connection to the database until this is successful
	var db *gorm.DB
	err := backoff.Retry(func() error {
		dbConnection, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")))
		if err != nil {
			return err
		}

		db = dbConnection

		return nil
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.Transaction{})
	if err != nil {
		panic(err)
	}

	return db
}
