package main

import (
	"estoque/config"

	"github.com/gin-gonic/gin"
)


func main() {
	cfg := config.Load()
	db := cfg.DBUrl
	var router *gin.Engine = gin.Default()
	
	
	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "API is running",
			"status": "success",
		})
	})
	
}