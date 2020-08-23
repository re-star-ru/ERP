package product

import (
	"backend/internal/app/models"
	"context"
)

type Usecase interface {
	CreateProduct(ctx context.Context, user *models.User, product models.Product) error // return created ?
	UpdateProduct(ctx context.Context, user *models.User, product models.Product) error // return updated ?

	GetOne(ctx context.Context, id string) (models.Product, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
	AddImage(ctx context.Context, user *models.User, productID string) error
}
