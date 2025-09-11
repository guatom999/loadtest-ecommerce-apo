package models

import "time"

// Persisted entities (DB layer structs)

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name" json:"firstName"`
	LastName     string    `db:"last_name" json:"lastName"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}

type Product struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	Price       float64    `db:"price" json:"price"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expiresAt,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
}
