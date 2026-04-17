package main

import (
	"estoque/config"
	"estoque/internal/handlers"
	"estoque/internal/repository"
	"estoque/internal/service"
	"github.com/gin-gonic/gin"
)


func main() {
	var router *gin.Engine = gin.Default()
	
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router.POST("/api/products", productHandler.AddProduct)
	router.GET("/api/products", productHandler.GetProducts)
	
	router.GET("/api", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "API is running",
			"status": "success",
		})
	})
	
	router.Run(":8080")
		
}