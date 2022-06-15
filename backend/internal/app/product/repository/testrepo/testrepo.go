package testrepo

import (
	"backend/internal/app/models"
	"backend/internal/app/product"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	db *mongo.Collection
}

func (r repository) GetByID(ctx context.Context, id string) (models.Product, error) {
	panic("implement me")
}

func (r repository) GetProducts(ctx context.Context) ([]models.Product, error) {
	panic("implement me")
}

func (r repository) CreateProduct(ctx context.Context, user *models.User, product models.Product) error {
	panic("implement me")
}

func (r repository) UpdateProduct(ctx context.Context, user *models.User, product models.Product) error {
	panic("implement me")
}

func NewRepository(collection string) product.Repository {
	return &repository{}
	//return &repository{DB.Collection(collection)}
}
