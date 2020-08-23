package productUsecase

import (
	"backend/internal/app/models"
	"backend/internal/app/product"
	"context"
	"time"
)

type usecase struct {
	productRepo    product.Repository
	contextTimeout time.Duration
}

func (p usecase) AddImage(ctx context.Context, user *models.User, productID string) error {

	panic("implement me")
}

//CreateProduct implemets ActionProductCreate
func (p usecase) CreateProduct(ctx context.Context, user *models.User, product models.Product) error {
	if !user.CheckPermission(models.ActionProductCreate) {
		return models.ErrAccessDenied
	}

	user.Sanitize()
	newProduct := models.Product{
		ID:              product.ID,
		Name:            product.Name,
		GUID:            product.GUID,
		SKU:             product.SKU,
		Description:     product.Description,
		Manufacturer:    product.Manufacturer,
		TypeGUID:        product.TypeGUID,
		TypeName:        product.TypeName,
		Characteristics: product.Characteristics,
		Properties:      product.Properties,
		Creator:         *user,
		CreatedAt:       time.Now(),
		LastModified:    time.Now(),
	}

	return p.productRepo.CreateProduct(ctx, user, newProduct)
}

func (p usecase) GetProducts(ctx context.Context) ([]models.Product, error) {
	return p.productRepo.GetProducts(ctx)
}

func (p usecase) GetOne(ctx context.Context, id string) (models.Product, error) {

	return p.productRepo.GetByID(ctx, id)
}

func NewUsecase(p product.Repository, timeout time.Duration) product.Usecase {
	return &usecase{
		p,
		timeout,
	}
}
