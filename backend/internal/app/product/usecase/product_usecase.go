package productUsecase

import (
	"backend/internal/app/domain"
	"context"
	"time"
)

type usecase struct {
	productRepo    domain.ProductRepository
	contextTimeout time.Duration
}

func NewUsecase(p domain.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &usecase{p, timeout}
}

func (p usecase) CreateProduct(ctx context.Context, user *domain.User, product *domain.Product) error {
	return p.productRepo.CreateProduct(ctx, user, product)
}

func (p usecase) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	return p.productRepo.GetProducts(ctx)
}
