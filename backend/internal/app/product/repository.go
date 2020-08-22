package product

import (
	"backend/internal/app/models"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (models.Product, error)
	GetProducts(ctx context.Context) ([]*models.Product, error)
	CreateProduct(ctx context.Context, user *models.User, product *models.Product) error
}
