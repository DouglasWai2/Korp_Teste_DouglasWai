package repository

import (
	"context"
	"database/sql"
	"estoque/internal/models"
	"fmt"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	AddProduct(ctx context.Context, product *models.Product) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(database *sql.DB) *productRepository {
	return &productRepository{db: database}
}

func (r *productRepository) AddProduct(ctx context.Context, product *models.Product) error {

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
			SELECT codigo, descricao, saldo
			FROM product
		`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Codigo, &product.Descricao, &product.Saldo); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
