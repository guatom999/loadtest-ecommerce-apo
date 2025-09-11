package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
)

type CartService interface {
	AddToCart(ctx context.Context, userID, productID string, qty int) error
	GetCart(ctx context.Context, userID string) (*models.CartResponse, error)
}

type cartService struct {
	rdb  *redis.Client
	prod repositories.ProductRepo
}

func NewCartService(rdb *redis.Client, prod repositories.ProductRepo) CartService {
	return &cartService{rdb: rdb, prod: prod}
}

func (s *cartService) cartKey(userID string) string { return "cart:" + userID }

func (s *cartService) AddToCart(ctx context.Context, userID, productID string, qty int) error {
	if qty <= 0 {
		return errors.New("quantity must be > 0")
	}
	// Ensure product exists
	if _, err := s.prod.Get(ctx, productID); err != nil {
		return err
	}

	key := s.cartKey(userID)
	if err := s.rdb.HIncrBy(ctx, key, productID, int64(qty)).Err(); err != nil {
		return err
	}
	// TTL 1 day (refresh on update)
	_ = s.rdb.Expire(ctx, key, 24*time.Hour).Err()
	return nil
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*models.CartResponse, error) {
	key := s.cartKey(userID)
	m, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	items := []models.CartItem{}
	for pid, qstr := range m {
		// parse quantity
		var qty int
		if _, err := fmt.Sscanf(qstr, "%d", &qty); err != nil {
			continue
		}
		p, err := s.prod.Get(ctx, pid)
		if err != nil {
			continue
		}
		items = append(items, models.CartItem{Product: *p, Quantity: qty})
	}
	return &models.CartResponse{Items: items}, nil
}
