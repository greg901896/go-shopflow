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

func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, id int64, name *string, price *string, stock *int) (*model.Product, error) {
	return s.repo.Update(ctx, id, name, price, stock)
}

func (s *ProductService) List(ctx context.Context, page, limit int) ([]model.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.List(ctx, limit, offset)
}
