package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	estoqueclient "faturamento/internal/clients/estoque"
	"faturamento/internal/models"
)

var (
	ErrNotaFiscalNotFound      = errors.New("nota fiscal not found")
	ErrNotaFiscalAlreadyClosed = errors.New("nota fiscal already closed")
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrProductNotFound         = errors.New("product not found")
)

type NotaFiscalRepository interface {
	GetNotasFiscais(ctx context.Context) ([]models.NotaFiscal, error)
	AddNotaFiscal(ctx context.Context, notaFiscal *models.NotaFiscal, itens []NotaFiscalItemInput) error
	PrintNotaFiscal(ctx context.Context, numero int64) (*models.NotaFiscal, error)
}

type notaFiscalRepository struct {
	db            *sql.DB
	estoqueClient *estoqueclient.Client
}

type notaFiscalItem struct {
	CodigoProduto string
	Quantidade    int
}

type NotaFiscalItemInput struct {
	CodigoProduto string
	Quantidade    int
}

func NewNotaFiscalRepository(database *sql.DB, estoqueClient *estoqueclient.Client) *notaFiscalRepository {
	return &notaFiscalRepository{
		db:            database,
		estoqueClient: estoqueClient,
	}
}

func (r *notaFiscalRepository) AddNotaFiscal(ctx context.Context, notaFiscal *models.NotaFiscal, itens []NotaFiscalItemInput) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
		INSERT INTO notas_fiscais (status)
		VALUES ($1)
		RETURNING numero
	`

	if err := tx.QueryRowContext(ctx, query, notaFiscal.Status).Scan(&notaFiscal.Numero); err != nil {
		return fmt.Errorf("add nota fiscal error: %w", err)
	}

	for _, item := range itens {
		if _, err := r.estoqueClient.GetProductByCode(ctx, item.CodigoProduto); err != nil {
			if errors.Is(err, estoqueclient.ErrProductNotFound) {
				return fmt.Errorf("%w: %s", ErrProductNotFound, item.CodigoProduto)
			}
			return fmt.Errorf("validate product in estoque error: %w", err)
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO nf_produtos (codigo_produto, numero_nf, quantidade)
			VALUES ($1, $2, $3)
		`, item.CodigoProduto, notaFiscal.Numero, item.Quantidade); err != nil {
			return fmt.Errorf("add nota fiscal items error: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	tx = nil

	return nil
}

func (r *notaFiscalRepository) GetNotasFiscais(ctx context.Context) ([]models.NotaFiscal, error) {
	var notasFiscais []models.NotaFiscal

	query := `
		SELECT numero, status
		FROM notas_fiscais
		ORDER BY numero
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notaFiscal models.NotaFiscal
		if err := rows.Scan(&notaFiscal.Numero, &notaFiscal.Status); err != nil {
			return nil, err
		}
		notasFiscais = append(notasFiscais, notaFiscal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notasFiscais, nil
}

func (r *notaFiscalRepository) PrintNotaFiscal(ctx context.Context, numero int64) (*models.NotaFiscal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	var notaFiscal models.NotaFiscal

	selectQuery := `
		SELECT numero, status
		FROM notas_fiscais
		WHERE numero = $1
		FOR UPDATE
	`

	if err := tx.QueryRowContext(ctx, selectQuery, numero).Scan(&notaFiscal.Numero, &notaFiscal.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotaFiscalNotFound
		}
		return nil, err
	}

	if notaFiscal.Status == "Fechada" {
		return nil, ErrNotaFiscalAlreadyClosed
	}

	items, err := r.getNotaFiscalItems(ctx, tx, numero)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if _, err := r.estoqueClient.DecrementStock(ctx, item.CodigoProduto, item.Quantidade); err != nil {
			switch {
			case errors.Is(err, estoqueclient.ErrProductNotFound):
				return nil, fmt.Errorf("%w: %s", ErrProductNotFound, item.CodigoProduto)
			case errors.Is(err, estoqueclient.ErrInsufficientStock):
				return nil, fmt.Errorf("%w for product %s", ErrInsufficientStock, item.CodigoProduto)
			default:
				return nil, fmt.Errorf("update product stock in estoque error: %w", err)
			}
		}
	}

	updateQuery := `
		UPDATE notas_fiscais
		SET status = 'Fechada'
		WHERE numero = $1
	`

	if _, err := tx.ExecContext(ctx, updateQuery, numero); err != nil {
		return nil, fmt.Errorf("print nota fiscal error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	tx = nil

	notaFiscal.Status = "Fechada"
	return &notaFiscal, nil
}

func (r *notaFiscalRepository) getNotaFiscalItems(ctx context.Context, tx *sql.Tx, numero int64) ([]notaFiscalItem, error) {
	query := `
		SELECT codigo_produto, SUM(quantidade)::BIGINT AS quantidade
		FROM nf_produtos
		WHERE numero_nf = $1
		GROUP BY codigo_produto
		ORDER BY codigo_produto
	`

	rows, err := tx.QueryContext(ctx, query, numero)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []notaFiscalItem
	for rows.Next() {
		var item notaFiscalItem
		if err := rows.Scan(&item.CodigoProduto, &item.Quantidade); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
