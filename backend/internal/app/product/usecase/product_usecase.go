package productUsecase

import (
	"backend/internal/app/models"
	"context"
	"time"
)

type usecase struct {
	productRepo    models.ProductRepository
	contextTimeout time.Duration
}

func (p usecase) AddImage(ctx context.Context, user *models.User, productID string) error {
	panic("implement me")
}

//CreateProduct implemets ActionProductCreate
func (p usecase) CreateProduct(ctx context.Context, user *models.User, product *models.Product) error {
	if !user.CheckPermission(models.ActionProductCreate) {
		return models.ErrAccessDenied
	}
	return p.productRepo.CreateProduct(ctx, user, product)
}

func (p usecase) GetProducts(ctx context.Context) ([]*models.Product, error) {
	return p.productRepo.GetProducts(ctx)
}

func (p usecase) GetOne(ctx context.Context, id string) (models.Product, error) {
	return models.Product{}, nil
}

func NewUsecase(p models.ProductRepository, timeout time.Duration) models.ProductUsecase {
	return &usecase{
		p,
		timeout,
	}
}
