package repository

import (
	"context"
	"database/sql"
	"errors"
	"estoque/internal/models"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductInUse         = errors.New("product in use")
	ErrProductNotFound      = errors.New("product not found")
	ErrInsufficientStock    = errors.New("insufficient stock")
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductByCode(ctx context.Context, codigo string) (*models.Product, error)
	AddProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, codigo string) error
	DecrementStock(ctx context.Context, codigo string, quantidade int) (*models.Product, error)
	IncrementStock(ctx context.Context, codigo string, quantidade int) (*models.Product, error)
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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrProductAlreadyExists
		}
		return fmt.Errorf("Add Product error: %w", err)
	}

	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, codigo string) error {
	query := `
		DELETE FROM product
		WHERE codigo = $1
	`

	result, err := r.db.ExecContext(ctx, query, codigo)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return ErrProductInUse
		}
		return fmt.Errorf("delete product error: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
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

func (r *productRepository) GetProductByCode(ctx context.Context, codigo string) (*models.Product, error) {
	var product models.Product

	query := `
		SELECT codigo, descricao, saldo
		FROM product
		WHERE codigo = $1
	`

	if err := r.db.QueryRowContext(ctx, query, codigo).Scan(&product.Codigo, &product.Descricao, &product.Saldo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) DecrementStock(ctx context.Context, codigo string, quantidade int) (*models.Product, error) {
	query := `
		UPDATE product
		SET saldo = saldo - $2
		WHERE codigo = $1
		  AND saldo >= $2
		RETURNING codigo, descricao, saldo
	`

	var product models.Product
	if err := r.db.QueryRowContext(ctx, query, codigo, quantidade).Scan(&product.Codigo, &product.Descricao, &product.Saldo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			exists, existsErr := r.productExists(ctx, codigo)
			if existsErr != nil {
				return nil, existsErr
			}
			if !exists {
				return nil, ErrProductNotFound
			}
			return nil, ErrInsufficientStock
		}
		return nil, fmt.Errorf("decrement stock error: %w", err)
	}

	return &product, nil
}

func (r *productRepository) IncrementStock(ctx context.Context, codigo string, quantidade int) (*models.Product, error) {
	query := `
		UPDATE product
		SET saldo = saldo + $2
		WHERE codigo = $1
		RETURNING codigo, descricao, saldo
	`

	var product models.Product
	if err := r.db.QueryRowContext(ctx, query, codigo, quantidade).Scan(&product.Codigo, &product.Descricao, &product.Saldo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("increment stock error: %w", err)
	}

	return &product, nil
}

func (r *productRepository) productExists(ctx context.Context, codigo string) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM product
			WHERE codigo = $1
		)
	`

	if err := r.db.QueryRowContext(ctx, query, codigo).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
