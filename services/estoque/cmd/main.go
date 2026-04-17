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

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
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

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	router.POST("/api/products", productHandler.AddProduct)
	router.GET("/api/products", productHandler.GetProducts)
	router.DELETE("/api/products/:codigo", productHandler.DeleteProduct)
	router.GET("/api/products/:codigo", productHandler.GetProductByCode)
	router.PATCH("/api/products/:codigo/increment", productHandler.IncrementStock)
	router.PATCH("/api/products/:codigo/decrement", productHandler.DecrementStock)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
			"status":  "success",
		})
	})

	router.Run(":8080")

}
