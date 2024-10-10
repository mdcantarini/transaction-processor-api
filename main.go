package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mdcantarini/transaction-processor-api/pkg/model"
	"github.com/mdcantarini/transaction-processor-api/pkg/service"
)

func main() {
	// set up the service
	s := service.NewService()

	// fill the database with mandatory initial data
	err := seedDB(s)
	if err != nil {
		panic(err)
	}

	// set up http router
	router := gin.Default()
	router.POST("/transactions/run-daily-report", s.RunDailyReport)

	err = router.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func seedDB(s *service.Service) error {
	accounts := []model.Account{
		{Email: "martin.d.cantarini@gmail.com"},
	}

	return s.AccountRepo().UpsertAccounts(accounts)
}
