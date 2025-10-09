package repositories

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/stretchr/testify/mock"
)

type (
	ProductRepoMock struct {
		mock.Mock
	}
)

func NewProductRepoMock() *ProductRepoMock {
	return &ProductRepoMock{}
}

func (m *ProductRepoMock) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepoMock) Get(ctx context.Context, id string) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepoMock) List(ctx context.Context, limit, offset int) ([]models.Product, int, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]models.Product), args.Int(1), args.Error(2)
}
