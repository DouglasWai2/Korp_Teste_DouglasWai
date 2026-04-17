package handlers

import(
	"github.com/gin-gonic/gin"
)

type ProductHandler interface{
	Service ProductService
}

func (p *ProductHandler) AddProduct(c* gin.Context){
	
	

    c.JSON(200, gin.H{
    	"status": "ok",
    })
}