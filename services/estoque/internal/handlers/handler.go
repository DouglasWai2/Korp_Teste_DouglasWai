package handlers

import (
	"net/http"

	"estoque/internal/models"
	"estoque/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (p *ProductHandler) AddProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if product.Codigo == "" || product.Descricao == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "codigo and descricao are required",
		})
		return
	}

	if err := p.ProductService.AddProduct(c.Request.Context(), product.Codigo, product.Descricao, product.Saldo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to add product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "product created",
		"data":    product,
	})
}

func (p *ProductHandler) GetProducts(c *gin.Context) {
	products, err := p.ProductService.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch products",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   products,
	})
}
