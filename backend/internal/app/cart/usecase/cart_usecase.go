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

func NewUsecase(c domain.CartRepository, time time.Duration) domain.CartUsecase {
	return &cart{c, time}
}

func (c cart) AddToCart(ctx context.Context, u *domain.User, productID string, count int) error {
	panic("implement me")
}

func (c cart) RemoveFromCart(ctx context.Context, u *domain.User, cartKey string) error {
	panic("implement me")
}
