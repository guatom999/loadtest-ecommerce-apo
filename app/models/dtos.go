package models

// Request/Response DTOs (handler <-> service)

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gte=0"`
	ExpiresAt   *string `json:"expiresAt"` // ISO date (YYYY-MM-DD), optional
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type AddToCartRequest struct {
	ProductID string `json:"productId" validate:"required,uuid4"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type CartItem struct {
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type CartResponse struct {
	Items []CartItem `json:"items"`
}
