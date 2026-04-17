package service

import (
	"context"

	"estoque/internal/models"
	"estoque/internal/repository"
)

type ProductService interface {
	AddProduct(ctx context.Context, codigo string, descricao string, saldo int) error
	GetProducts(ctx context.Context) ([]models.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (p *productService) AddProduct(ctx context.Context, codigo string, descricao string, saldo int) error {
	product := models.Product{
		Codigo:    codigo,
		Descricao: descricao,
		Saldo:     saldo,
	}

	err := p.repo.AddProduct(ctx, &product)
	return err
}

func (p *productService) GetProducts(ctx context.Context) ([]models.Product, error) {
	products, err := p.repo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
