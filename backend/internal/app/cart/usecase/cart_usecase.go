package cartUsecase

import (
	"backend/internal/app/cart"
	"backend/internal/app/models"
	"context"
	"time"
)

type usecase struct {
	cartRepo       cart.Repository
	contextTimeout time.Duration
}

func (c usecase) ShowUsersCart(ctx context.Context, u *models.User) (*models.Cart, error) {
	return c.cartRepo.GetUsersCart(ctx, u)
}

func (c usecase) AddProductToCart(ctx context.Context, u *models.User, productID string, count int) error {
	return c.cartRepo.AddToCart(ctx, u, productID, count)
}

func (c usecase) RemoveProductFromCart(ctx context.Context, u *models.User, cartKey string) error {
	panic("implement me")
}

func NewUsecase(c cart.Repository, time time.Duration) cart.Usecase {
	return &usecase{c, time}
}
