package repository

import (
	"context"
	"database/sql"
	"estoque/internal/models"
	"fmt"
)


type ProductRepository interface {
		GetProducts()[]models.Product
		AddProduct(product *models.Product) error
}

type productRepository struct{
	db *sql.DB
}

func NewProductRepository(database *sql.DB) *productRepository{
	return &productRepository{db: database}
}

func (r *productRepository) AddProduct(ctx context.Context, product models.Product) error {
	
	query := `
			INSERT INTO product (codigo, descricao, saldo)
			VALUES ($1, $2, $3)
		`
	_, err := r.db.ExecContext(ctx, query, product.Codigo, product.Descricao, product.Saldo)
	
	if err != nil {
		return fmt.Errorf("Add Product error: %w", err)
	}
	
	return nil
}

func (r *productRepository) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	query := `
			SELECT *
			FROM product
		`
	if err := r.db.QueryRowContext(ctx, query).Scan(&products); err != nil {
		return nil, fmt.Errorf("Get Products error: %w", err)
	}


	return products, nil
}