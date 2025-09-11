package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
	"github.com/guatom999/ecommerce-product-api/app/utils"
)

type AuthService interface {
	Register(ctx context.Context, in models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, in models.LoginRequest) (string, error)
}

type authService struct {
	users repositories.UserRepo
	jwt   *utils.JWTMaker
}

func NewAuthService(users repositories.UserRepo, jwt *utils.JWTMaker) AuthService {
	return &authService{users: users, jwt: jwt}
}

func (s *authService) Register(ctx context.Context, in models.RegisterRequest) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{Email: in.Email, PasswordHash: string(hash), FirstName: in.FirstName, LastName: in.LastName}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *authService) Login(ctx context.Context, in models.LoginRequest) (string, error) {
	u, err := s.users.FindByEmail(ctx, in.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		return "", errors.New("invalid credentials")
	}
	return s.jwt.Create(u.ID)
}
