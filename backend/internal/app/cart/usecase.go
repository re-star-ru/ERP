package cart

import (
	"backend/internal/app/models"
	"context"
)

type Usecase interface {
	AddProductToCart(ctx context.Context, u *models.User, productID string, count int) error
	RemoveProductFromCart(ctx context.Context, u *models.User, cartKey string) error
	ShowUsersCart(ctx context.Context, u *models.User) (*models.Cart, error)
}
