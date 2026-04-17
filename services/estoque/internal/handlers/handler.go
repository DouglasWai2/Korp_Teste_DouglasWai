package handlers

import (
	"errors"
	"net/http"

	"estoque/internal/models"
	"estoque/internal/repository"
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
		switch {
		case errors.Is(err, repository.ErrProductAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "product already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to add product",
				"error":   err.Error(),
			})
		}
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

func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	err := p.ProductService.DeleteProduct(c.Request.Context(), c.Param("codigo"))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "product not found",
			})
		case errors.Is(err, repository.ErrProductInUse):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "product is linked to existing invoices",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to delete product",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "product deleted",
	})
}

func (p *ProductHandler) GetProductByCode(c *gin.Context) {
	product, err := p.ProductService.GetProductByCode(c.Request.Context(), c.Param("codigo"))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "product not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to fetch product",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   product,
	})
}

func (p *ProductHandler) DecrementStock(c *gin.Context) {
	var request models.StockChangeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if request.Quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "quantidade must be greater than zero",
		})
		return
	}

	product, err := p.ProductService.DecrementStock(c.Request.Context(), c.Param("codigo"), request.Quantidade)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "product not found",
			})
		case errors.Is(err, repository.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "insufficient stock",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to decrement stock",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "stock updated",
		"data":    product,
	})
}

func (p *ProductHandler) IncrementStock(c *gin.Context) {
	var request models.StockChangeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if request.Quantidade <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "quantidade must be greater than zero",
		})
		return
	}

	product, err := p.ProductService.IncrementStock(c.Request.Context(), c.Param("codigo"), request.Quantidade)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "product not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to increment stock",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "stock updated",
		"data":    product,
	})
}
