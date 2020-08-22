package product

import (
	"backend/internal/app/models"
	"context"
)

type usecase interface {
	CreateProduct(ctx context.Context, user *models.Characteristic, product *models.Product) error
	GetProducts(ctx context.Context) ([]*models.Product, error)
	AddImage(ctx context.Context, user *models.User, productID string) error
	GetOne(ctx context.Context, id string) (models.Product, error)
}
