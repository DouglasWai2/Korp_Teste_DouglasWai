package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"faturamento/internal/models"
)

var (
	ErrNotaFiscalNotFound      = errors.New("nota fiscal not found")
	ErrNotaFiscalAlreadyClosed = errors.New("nota fiscal already closed")
)

type NotaFiscalRepository interface {
	GetNotasFiscais(ctx context.Context) ([]models.NotaFiscal, error)
	AddNotaFiscal(ctx context.Context, notaFiscal *models.NotaFiscal) error
	PrintNotaFiscal(ctx context.Context, numero int64) (*models.NotaFiscal, error)
}

type notaFiscalRepository struct {
	db *sql.DB
}

func NewNotaFiscalRepository(database *sql.DB) *notaFiscalRepository {
	return &notaFiscalRepository{db: database}
}

func (r *notaFiscalRepository) AddNotaFiscal(ctx context.Context, notaFiscal *models.NotaFiscal) error {
	query := `
		INSERT INTO notas_fiscais (status)
		VALUES ($1)
		RETURNING numero
	`

	if err := r.db.QueryRowContext(ctx, query, notaFiscal.Status).Scan(&notaFiscal.Numero); err != nil {
		return fmt.Errorf("add nota fiscal error: %w", err)
	}

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
	var notaFiscal models.NotaFiscal

	selectQuery := `
		SELECT numero, status
		FROM notas_fiscais
		WHERE numero = $1
	`

	if err := r.db.QueryRowContext(ctx, selectQuery, numero).Scan(&notaFiscal.Numero, &notaFiscal.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotaFiscalNotFound
		}
		return nil, err
	}

	if notaFiscal.Status == "Fechada" {
		return nil, ErrNotaFiscalAlreadyClosed
	}

	updateQuery := `
		UPDATE notas_fiscais
		SET status = 'Fechada'
		WHERE numero = $1
	`

	if _, err := r.db.ExecContext(ctx, updateQuery, numero); err != nil {
		return nil, fmt.Errorf("print nota fiscal error: %w", err)
	}

	notaFiscal.Status = "Fechada"
	return &notaFiscal, nil
}
