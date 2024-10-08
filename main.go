package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mdcantarini/transaction-processor-api/pkg/model"
)

func main() {
	// set up the service
	service := NewService()

	// fill the database with mandatory initial data
	err := fillDB(service)
	if err != nil {
		panic(err)
	}

	// set up http router
	router := gin.Default()
	SetRoutes(router, service)
	err = router.Run(":8000")
	if err != nil {
		panic(err)
	}
}

func fillDB(s *Service) error {
	accounts := []model.Account{
		{Email: "martin.d.cantarini@gmail.com"},
	}

	return s.accountRepo.UpsertAccounts(accounts)
}
