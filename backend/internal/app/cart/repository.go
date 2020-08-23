package cart

import (
	"backend/internal/app/models"
	"context"
)

type Repository interface {
	AddToCart(ctx context.Context, u *models.User, productID string, count int) error
	RemoveFromCart(ctx context.Context, u *models.User, cartKey string) error
	GetUsersCart(ctx context.Context, u *models.User) (*models.Cart, error)
}
