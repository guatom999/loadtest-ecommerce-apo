package repositories

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/stretchr/testify/mock"
)

type (
	UserRepoMock struct {
		mock.Mock
	}
)

func NewUserRepoMock() UserRepo {
	return &UserRepoMock{}
}

func (m *UserRepoMock) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepoMock) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
