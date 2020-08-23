package auth

import (
	"backend/internal/app/models"
	"context"
)

// UserRepository represents the user's repository contract
type Repository interface {
	GetByID(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user *models.User) error
	ValidateUser(ctx context.Context, user *models.User) (bool, error)
}
