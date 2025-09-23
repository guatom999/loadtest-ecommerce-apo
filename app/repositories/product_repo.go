package repositories

import (
	"context"
	"log"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/jmoiron/sqlx"
)

type ProductRepo interface {
	Create(ctx context.Context, p *models.Product) error
	Get(ctx context.Context, id string) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]models.Product, int, error)
}

type productRepo struct{ db *sqlx.DB }

func NewProductRepo(db *sqlx.DB) ProductRepo { return &productRepo{db: db} }

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	q := `INSERT INTO products (id, name, description, price, expires_at) VALUES (uuid_generate_v4(), $1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, q, p.Name, p.Description, p.Price, p.ExpiresAt).Scan(&p.ID, &p.CreatedAt)
}

func (r *productRepo) Get(ctx context.Context, id string) (*models.Product, error) {
	var p models.Product
	err := r.db.GetContext(ctx, &p, `SELECT * FROM products WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) List(ctx context.Context, limit, offset int) ([]models.Product, int, error) {
	list := []models.Product{}
	if err := r.db.SelectContext(ctx, &list, `SELECT * FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset); err != nil {
		log.Printf("failed to list products: %v", err)
		return nil, 0, err
	}
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM products`); err != nil {
		log.Printf("failed to count products: %v", err)
		return nil, 0, err
	}
	return list, total, nil
}
