package service

import (
	"context"

	"faturamento/internal/models"
	"faturamento/internal/repository"
)

type NotaFiscalService interface {
	AddNotaFiscal(ctx context.Context, status string) (*models.NotaFiscal, error)
	GetNotasFiscais(ctx context.Context) ([]models.NotaFiscal, error)
	PrintNotaFiscal(ctx context.Context, numero int64) (*models.NotaFiscal, error)
}

type notaFiscalService struct {
	repo repository.NotaFiscalRepository
}

func NewNotaFiscalService(repo repository.NotaFiscalRepository) *notaFiscalService {
	return &notaFiscalService{
		repo: repo,
	}
}

func (s *notaFiscalService) AddNotaFiscal(ctx context.Context, status string) (*models.NotaFiscal, error) {
	notaFiscal := models.NotaFiscal{
		Status: status,
	}

	if err := s.repo.AddNotaFiscal(ctx, &notaFiscal); err != nil {
		return nil, err
	}

	return &notaFiscal, nil
}

func (s *notaFiscalService) GetNotasFiscais(ctx context.Context) ([]models.NotaFiscal, error) {
	notasFiscais, err := s.repo.GetNotasFiscais(ctx)
	if err != nil {
		return nil, err
	}

	return notasFiscais, nil
}

func (s *notaFiscalService) PrintNotaFiscal(ctx context.Context, numero int64) (*models.NotaFiscal, error) {
	return s.repo.PrintNotaFiscal(ctx, numero)
}
