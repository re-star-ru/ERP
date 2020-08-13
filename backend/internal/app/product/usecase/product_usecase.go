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

func (p usecase) AddImage(ctx context.Context, user *domain.User, productID string) error {
	panic("implement me")
}

//CreateProduct implemets ActionProductCreate
func (p usecase) CreateProduct(ctx context.Context, user *domain.User, product *domain.Product) error {
	if !user.CheckPermission(domain.ActionProductCreate) {
		return domain.ErrAccessDenied
	}

	return p.productRepo.CreateProduct(ctx, user, product)
}

func (p usecase) GetProducts(ctx context.Context) ([]*domain.Product, error) {
	return p.productRepo.GetProducts(ctx)
}

func NewUsecase(p domain.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &usecase{p, timeout}
}
