package repositories

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	Create(ctx context.Context, u *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct{ db *sqlx.DB }

func NewUserRepo(db *sqlx.DB) UserRepo { return &userRepo{db: db} }

func (r *userRepo) Create(ctx context.Context, u *models.User) error {
	q := `INSERT INTO users (id, email, password_hash, first_name, last_name) VALUES (uuid_generate_v4(), $1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRowxContext(ctx, q, u.Email, u.PasswordHash, u.FirstName, u.LastName).Scan(&u.ID, &u.CreatedAt)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
