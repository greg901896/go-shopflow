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

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at
		FROM products
		WHERE id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int64, name *string, price *string, stock *int) (*model.Product, error) {
	query := `
		UPDATE products
		SET
			name  = COALESCE($2, name),
			price = COALESCE($3, price),
			stock = COALESCE($4, stock)
		WHERE id = $1
		RETURNING id, name, price, stock, created_at
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id, name, price, stock).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) List(ctx context.Context, limit, offset int) ([]model.Product, error) {
	query := `
		SELECT id, name, price, stock, created_at
		FROM products
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []model.Product{}
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}
