package main

import (
	"nf-system/api/estoque/config"

	"github.com/gin-gonic/gin"
)


func main() {
	db := config.Db_connect()
	var router *gin.Engine = gin.Default()
	
	
	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "API is running",
			"status": "success",
		})
	})
	
}