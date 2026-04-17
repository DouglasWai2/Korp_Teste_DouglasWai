package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"faturamento/internal/repository"
	"faturamento/internal/service"

	"github.com/gin-gonic/gin"
)

type NotaFiscalHandler struct {
	NotaFiscalService service.NotaFiscalService
}

type addNotaFiscalItemRequest struct {
	CodigoProduto string `json:"codigo_produto"`
	Quantidade    int    `json:"quantidade"`
}

type addNotaFiscalRequest struct {
	Status string                     `json:"status"`
	Itens  []addNotaFiscalItemRequest `json:"itens"`
}

func NewNotaFiscalHandler(notaFiscalService service.NotaFiscalService) *NotaFiscalHandler {
	return &NotaFiscalHandler{
		NotaFiscalService: notaFiscalService,
	}
}

func (h *NotaFiscalHandler) AddNotaFiscal(c *gin.Context) {
	var request addNotaFiscalRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if request.Status == "" {
		request.Status = "Aberta"
	}

	if request.Status != "Aberta" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "new nota fiscal must start with status 'Aberta'",
		})
		return
	}

	if len(request.Itens) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "itens is required",
		})
		return
	}

	itens := make([]service.AddNotaFiscalItemInput, 0, len(request.Itens))
	codigos := make(map[string]struct{}, len(request.Itens))
	for _, item := range request.Itens {
		if item.CodigoProduto == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "codigo_produto is required",
			})
			return
		}

		if item.Quantidade <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "quantidade must be greater than zero",
			})
			return
		}

		if _, exists := codigos[item.CodigoProduto]; exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "codigo_produto must be unique within itens",
			})
			return
		}
		codigos[item.CodigoProduto] = struct{}{}

		itens = append(itens, service.AddNotaFiscalItemInput{
			CodigoProduto: item.CodigoProduto,
			Quantidade:    item.Quantidade,
		})
	}

	notaFiscal, err := h.NotaFiscalService.AddNotaFiscal(c.Request.Context(), request.Status, itens)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to add nota fiscal",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "nota fiscal created",
		"data":    notaFiscal,
	})
}

func (h *NotaFiscalHandler) GetNotasFiscais(c *gin.Context) {
	notasFiscais, err := h.NotaFiscalService.GetNotasFiscais(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to fetch notas fiscais",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   notasFiscais,
	})
}

func (h *NotaFiscalHandler) PrintNotaFiscal(c *gin.Context) {
	numero, err := strconv.ParseInt(c.Param("numero"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid nota fiscal number",
		})
		return
	}

	notaFiscal, err := h.NotaFiscalService.PrintNotaFiscal(c.Request.Context(), numero)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotaFiscalNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "nota fiscal not found",
			})
		case errors.Is(err, repository.ErrNotaFiscalAlreadyClosed):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "nota fiscal is already closed",
			})
		case errors.Is(err, repository.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case errors.Is(err, repository.ErrInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "failed to print nota fiscal",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "nota fiscal printed",
		"data":    notaFiscal,
	})
}
