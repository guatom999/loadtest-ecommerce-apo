package services

import (
	"context"
	"time"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
)

type ProductService interface {
	Create(ctx context.Context, in models.CreateProductRequest) (*models.Product, error)
	Get(ctx context.Context, id string) (*models.Product, error)
	List(ctx context.Context, limit, offset int) ([]models.Product, int, error)
}

type productService struct{ repo repositories.ProductRepo }

func NewProductService(repo repositories.ProductRepo) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, in models.CreateProductRequest) (*models.Product, error) {
	var exp *time.Time
	if in.ExpiresAt != nil && *in.ExpiresAt != "" {
		if t, err := time.Parse("2006-01-02", *in.ExpiresAt); err == nil {
			exp = &t
		}
	}
	p := &models.Product{Name: in.Name, Description: in.Description, Price: in.Price, ExpiresAt: exp}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) Get(ctx context.Context, id string) (*models.Product, error) {
	return s.repo.Get(ctx, id)
}

func (s *productService) List(ctx context.Context, limit, offset int) ([]models.Product, int, error) {
	return s.repo.List(ctx, limit, offset)
}
