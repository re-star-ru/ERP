package models

import (
	"time"
)

// Product это продукт
type Product struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	GUID         string `json:"guid" storm:"id"`
	SKU          string `json:"sku" storm:"index"`
	Description  string `json:"description"`
	Manufacturer string `json:"manufacturer"`
	TypeGUID     string `json:"typeGUID"`
	TypeName     string `json:"typeName"`

	Characteristics []Characteristic `json:"characteristics"`

	Properties []Property `json:"properties"`

	Creator User `json:"creator"`

	CreatedAt    time.Time `json:"createdAt"`
	LastModified time.Time `json:"lastModified"`
}

//type ProductUsecase interface {
//	CreateProduct(ctx context.Context, user *User, product *Product) error
//	GetProducts(ctx context.Context) ([]*Product, error)
//	AddImage(ctx context.Context, user *User, productID string) error
//	GetOne(ctx context.Context, id string) (Product, error)
//}
//
//type ProductRepository interface {
//	GetByID(ctx context.Context, id string) (Product, error)
//	GetProducts(ctx context.Context) ([]*Product, error)
//	CreateProduct(ctx context.Context, user *User, product *Product) error
//}

//// Group это структура с артикулами
//type Group struct {
//	ID  int    `storm:"increment"`
//	SKU string `json:"sku" storm:"index"`
//
//	MainProductName string `json:"mainProductName"`
//	MainProductGUID string `json:"mainProductGUID"`
//
//	MainProductTypeName string `json:"mainProductTypeName"`
//	MainProductTypeGUID string `json:"mainProductTypeGUID"`
//
//	MainProductProperties []Property `json:"mainProductProperties"`
//
//	Products []Product `json:"catalog"`
//	//
//	//SkuConfig struct {
//	//} `json:"skuConfig"`
//}

type Characteristic struct {
	GUID  string `json:"characteristic" storm:"id"`
	Name  string `json:"characteristicName" storm:"index"`
	Owner string `json:"characteristicOwner"`
}

type Property struct {
	Name  string `json:"propertyName" storm:"index"`
	GUID  string `json:"property"`
	Value string `json:"value"`
	Unit  string `json:"unit"`
}
