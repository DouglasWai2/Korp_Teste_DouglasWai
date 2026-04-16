package handlers

import(
	"github.com/gin-gonic/gin"
)

func AddProduct(c* gin.Context){

    c.JSON(200, gin.H{
    	"status": "ok",
    })
}