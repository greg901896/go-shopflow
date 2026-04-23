package repository

import (
	"context"

	"github.com/greg901896/go-shopflow/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) error {
	query := `
		INSERT INTO products (name, price, stock)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, p.Name, p.Price, p.Stock).
		Scan(&p.ID, &p.CreatedAt)
}
