package service

import (
	"context"

	"faturamento/internal/models"
	"faturamento/internal/repository"
)

type AddNotaFiscalItemInput struct {
	CodigoProduto string
	Quantidade    int
}

type NotaFiscalService interface {
	AddNotaFiscal(ctx context.Context, status string, itens []AddNotaFiscalItemInput) (*models.NotaFiscal, error)
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

func (s *notaFiscalService) AddNotaFiscal(ctx context.Context, status string, itens []AddNotaFiscalItemInput) (*models.NotaFiscal, error) {
	notaFiscal := models.NotaFiscal{
		Status: status,
	}

	repoItems := make([]repository.NotaFiscalItemInput, 0, len(itens))
	for _, item := range itens {
		repoItems = append(repoItems, repository.NotaFiscalItemInput{
			CodigoProduto: item.CodigoProduto,
			Quantidade:    item.Quantidade,
		})
	}

	if err := s.repo.AddNotaFiscal(ctx, &notaFiscal, repoItems); err != nil {
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
