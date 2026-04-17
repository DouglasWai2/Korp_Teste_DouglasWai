package main

import (
	"faturamento/config"
	"faturamento/internal/handlers"
	"faturamento/internal/repository"
	"faturamento/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	var router *gin.Engine = gin.Default()

	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	notaFiscalRepo := repository.NewNotaFiscalRepository(db)
	notaFiscalService := service.NewNotaFiscalService(notaFiscalRepo)
	notaFiscalHandler := handlers.NewNotaFiscalHandler(notaFiscalService)

	router.POST("/api/notas-fiscais", notaFiscalHandler.AddNotaFiscal)
	router.GET("/api/notas-fiscais", notaFiscalHandler.GetNotasFiscais)
	router.PATCH("/api/notas-fiscais/:numero/print", notaFiscalHandler.PrintNotaFiscal)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
			"status":  "success",
		})
	})

	router.Run(":8081")
}
