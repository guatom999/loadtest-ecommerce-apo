package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/guatom999/ecommerce-product-api/app/models"
	"github.com/guatom999/ecommerce-product-api/app/repositories"
	"github.com/guatom999/ecommerce-product-api/app/utils"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		input       models.RegisterRequest
		setupMock   func(*repositories.UserRepoMock)
		expectUser  bool
		expectedErr string
	}{
		{
			name: "successful_registration",
			input: models.RegisterRequest{
				Email:     "newuser@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				// Check if user exists (should return not found)
				mockRepo.On("GetByEmail", mock.Anything, "newuser@example.com").Return(nil, errors.New("user not found"))
				// Create user should succeed
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					user := args.Get(1).(*models.User)
					user.ID = "new-user-id"
				})
			},
			expectUser: true,
		},
		{
			name: "email_already_exists",
			input: models.RegisterRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "Jane",
				LastName:  "Doe",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				existingUser := &models.User{
					ID:        "existing-id",
					Email:     "existing@example.com",
					FirstName: "Jane",
					LastName:  "Doe",
				}
				mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedErr: "email already in use",
		},
		{
			name: "database_error_on_create",
			input: models.RegisterRequest{
				Email:     "newuser2@example.com",
				Password:  "password123",
				FirstName: "Bob",
				LastName:  "Smith",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				mockRepo.On("GetByEmail", mock.Anything, "newuser2@example.com").Return(nil, errors.New("user not found"))
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error"))
			},
			expectedErr: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(repositories.UserRepoMock)
			jwt := utils.NewJWTMaker()
			tt.setupMock(mockRepo)
			service := NewAuthService(mockRepo, jwt)
			ctx := context.Background()

			// Execute
			user, err := service.Register(ctx, tt.input)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, user)
			} else if tt.expectUser {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.input.Email, user.Email)
				assert.Equal(t, tt.input.FirstName, user.FirstName)
				assert.Equal(t, tt.input.LastName, user.LastName)
				assert.NotEmpty(t, user.ID)
				assert.Empty(t, user.PasswordHash) // Password hash should not be in response
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Pre-hash a password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		input       models.LoginRequest
		setupMock   func(*repositories.UserRepoMock)
		expectToken bool
		expectedErr string
	}{
		{
			name: "successful_login",
			input: models.LoginRequest{
				Email:    "user@example.com",
				Password: "correctpassword",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				user := &models.User{
					ID:           "user-id-123",
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
				}
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(user, nil)
			},
			expectToken: true,
		},
		{
			name: "user_not_found",
			input: models.LoginRequest{
				Email:    "notfound@example.com",
				Password: "anypassword",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("user not found"))
			},
			expectedErr: "invalid credentials",
		},
		{
			name: "wrong_password",
			input: models.LoginRequest{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				user := &models.User{
					ID:           "user-id-123",
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
					FirstName:    "Test",
					LastName:     "User",
				}
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(user, nil)
			},
			expectedErr: "invalid credentials",
		},
		{
			name: "database_error",
			input: models.LoginRequest{
				Email:    "user@example.com",
				Password: "anypassword",
			},
			setupMock: func(mockRepo *repositories.UserRepoMock) {
				mockRepo.On("GetByEmail", mock.Anything, "user@example.com").Return(nil, errors.New("database connection failed"))
			},
			expectedErr: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(repositories.UserRepoMock)
			tt.setupMock(mockRepo)
			jwt := utils.NewJWTMaker()
			service := NewAuthService(mockRepo, jwt)
			ctx := context.Background()

			// Execute
			token, err := service.Login(ctx, tt.input)

			// Assert
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Empty(t, token)
			} else if tt.expectToken {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				// Token should be a valid JWT (basic check)
				assert.Contains(t, token, ".")
				assert.True(t, len(token) > 20)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
