package service

import (
	"github.com/gin-gonic/gin"

	"estoque/internal/models"
	"estoque/internal/repository"
)

type ProductService interface{
	AddProduct(codigo string, descricao string, saldo number)(*Product, error)
}

func AddProduct(c *gin.Context, codigo string, descricao string, saldo number){
	product := models.Product{
		Codigo:    codigo,
		Descricao: descricao,
		Saldo:     saldo,
	}
	repository.CreateProduct(product)
}