package main

import (
	"faturamento/config"
	estoqueclient "faturamento/internal/clients/estoque"
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

	estoqueClient := estoqueclient.NewClient(config.GetEstoqueAPIURL())
	notaFiscalRepo := repository.NewNotaFiscalRepository(db, estoqueClient)
	notaFiscalService := service.NewNotaFiscalService(notaFiscalRepo)
	notaFiscalHandler := handlers.NewNotaFiscalHandler(notaFiscalService)

	router.GET("/api/faturamento", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
			"status":  "success",
		})
	})

	faturamento := router.Group("/api/faturamento")
	{
		faturamento.POST("/notas-fiscais", notaFiscalHandler.AddNotaFiscal)
		faturamento.GET("/notas-fiscais", notaFiscalHandler.GetNotasFiscais)
		faturamento.PATCH("/notas-fiscais/:numero/imprimir", notaFiscalHandler.PrintNotaFiscal)
	}

	router.Run(":8081")
}
