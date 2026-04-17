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

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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
		faturamento.GET("/notas-fiscais/:numero", notaFiscalHandler.GetNotaFiscalByNumero)
		faturamento.PATCH("/notas-fiscais/:numero/imprimir", notaFiscalHandler.PrintNotaFiscal)
	}

	router.Run(":8081")
}
