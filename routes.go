package main

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(router *gin.Engine, s *Service) {
	router.POST("/transactions/run-daily-report", s.RunDailyReport)
}
