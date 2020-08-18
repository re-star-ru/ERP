package cartUsecase

import (
	"backend/internal/app/domain"
	"context"
	"time"
)

type cart struct {
	cartRepo       domain.CartRepository
	contextTimeout time.Duration
}

func (c cart) ShowUsersCart(ctx context.Context, u *domain.User) (*domain.Cart, error) {
	return c.cartRepo.GetUsersCart(ctx, u)
}

func (c cart) AddProductToCart(ctx context.Context, u *domain.User, productID string, count int) error {
	return c.cartRepo.AddToCart(ctx, u, productID, count)
}

func (c cart) RemoveProductFromCart(ctx context.Context, u *domain.User, cartKey string) error {
	panic("implement me")
}

func NewUsecase(c domain.CartRepository, time time.Duration) domain.CartUsecase {
	return &cart{c, time}
}
