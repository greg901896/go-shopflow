package service

import (
	"context"

	"github.com/greg901896/go-shopflow/internal/model"
	"github.com/greg901896/go-shopflow/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	return s.repo.Create(ctx, p)
}
