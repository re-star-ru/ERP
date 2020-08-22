package cartUsecase

import (
	"backend/internal/app/models"
	"context"
	"time"
)

type cart struct {
	cartRepo       models.CartRepository
	contextTimeout time.Duration
}

func (c cart) ShowUsersCart(ctx context.Context, u *models.User) (*models.Cart, error) {
	return c.cartRepo.GetUsersCart(ctx, u)
}

func (c cart) AddProductToCart(ctx context.Context, u *models.User, productID string, count int) error {
	return c.cartRepo.AddToCart(ctx, u, productID, count)
}

func (c cart) RemoveProductFromCart(ctx context.Context, u *models.User, cartKey string) error {
	panic("implement me")
}

func NewUsecase(c models.CartRepository, time time.Duration) models.CartUsecase {
	return &cart{c, time}
}
