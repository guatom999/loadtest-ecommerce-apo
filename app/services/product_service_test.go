package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
)

func TestProductService_Create(t *testing.T) {
	tests := []struct {
		name        string
		input       models.CreateProductRequest
		setupMock   func(*repositories.ProductRepoMock)
		expected    *models.Product
		expectedErr string
	}{
		{
			name: "successful_create_product",
			input: models.CreateProductRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       99.99,
			},
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(1).(*models.Product)
					product.ID = "test-id-123"
					product.CreatedAt = time.Now()
				})
			},
			expected: &models.Product{
				ID:          "test-id-123",
				Name:        "Test Product",
				Description: "Test Description",
				Price:       99.99,
			},
		},
		{
			name: "create_product_with_expiry_date",
			input: models.CreateProductRequest{
				Name:        "Expiring Product",
				Description: "Will expire soon",
				Price:       49.99,
				ExpiresAt:   stringPtr("2024-12-31"),
			},
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					product := args.Get(1).(*models.Product)
					product.ID = "expiring-id-456"
					product.CreatedAt = time.Now()
				})
			},
			expected: &models.Product{
				ID:          "expiring-id-456",
				Name:        "Expiring Product",
				Description: "Will expire soon",
				Price:       49.99,
			},
		},
		{
			name: "repository_error",
			input: models.CreateProductRequest{
				Name:        "Failed Product",
				Description: "Should fail",
				Price:       1.00,
			},
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Product")).Return(errors.New("database error"))
			},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(repositories.ProductRepoMock)
			tt.setupMock(mockRepo)
			service := NewProductService(mockRepo)
			ctx := context.Background()

			// Execute
			result, err := service.Create(ctx, tt.input)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Description, result.Description)
				assert.Equal(t, tt.expected.Price, result.Price)
				assert.NotEmpty(t, result.ID)
				assert.NotZero(t, result.CreatedAt)

				// Check expiry date if provided
				if tt.input.ExpiresAt != nil {
					assert.NotNil(t, result.ExpiresAt)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_Get(t *testing.T) {
	tests := []struct {
		name        string
		productID   string
		setupMock   func(*repositories.ProductRepoMock)
		expected    *models.Product
		expectedErr string
	}{
		{
			name:      "successful_get_product",
			productID: "existing-id",
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				product := &models.Product{
					ID:          "existing-id",
					Name:        "Existing Product",
					Description: "Product exists",
					Price:       75.50,
					CreatedAt:   time.Now(),
				}
				mockRepo.On("Get", mock.Anything, "existing-id").Return(product, nil)
			},
			expected: &models.Product{
				ID:          "existing-id",
				Name:        "Existing Product",
				Description: "Product exists",
				Price:       75.50,
			},
		},
		{
			name:      "product_not_found",
			productID: "non-existing-id",
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("Get", mock.Anything, "non-existing-id").Return((*models.Product)(nil), errors.New("product not found"))
			},
			expectedErr: "product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(repositories.ProductRepoMock)
			tt.setupMock(mockRepo)
			service := NewProductService(mockRepo)
			ctx := context.Background()

			// Execute
			result, err := service.Get(ctx, tt.productID)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Description, result.Description)
				assert.Equal(t, tt.expected.Price, result.Price)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProductService_List(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		offset        int
		setupMock     func(*repositories.ProductRepoMock)
		expectedCount int
		expectedTotal int
		expectedErr   string
	}{
		{
			name:   "successful_list_products",
			limit:  10,
			offset: 0,
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				products := []models.Product{
					{ID: "1", Name: "Product 1", Price: 10.00, CreatedAt: time.Now()},
					{ID: "2", Name: "Product 2", Price: 20.00, CreatedAt: time.Now()},
					{ID: "3", Name: "Product 3", Price: 30.00, CreatedAt: time.Now()},
				}
				mockRepo.On("List", mock.Anything, 10, 0).Return(products, 3, nil)
			},
			expectedCount: 3,
			expectedTotal: 3,
		},
		{
			name:   "empty_result",
			limit:  10,
			offset: 100,
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("List", mock.Anything, 10, 100).Return([]models.Product{}, 0, nil)
			},
			expectedCount: 0,
			expectedTotal: 0,
		},
		{
			name:   "repository_error",
			limit:  10,
			offset: 0,
			setupMock: func(mockRepo *repositories.ProductRepoMock) {
				mockRepo.On("List", mock.Anything, 10, 0).Return([]models.Product{}, 0, errors.New("database connection failed"))
			},
			expectedErr: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(repositories.ProductRepoMock)
			tt.setupMock(mockRepo)
			service := NewProductService(mockRepo)
			ctx := context.Background()

			// Execute
			products, total, err := service.List(ctx, tt.limit, tt.offset)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, products)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Len(t, products, tt.expectedCount)
				assert.Equal(t, tt.expectedTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
