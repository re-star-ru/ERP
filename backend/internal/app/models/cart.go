package models

import "context"

type Cart struct {
	OwnerID       string `json:"userID"`
	AddedProducts `json:"addedProducts"`
}

type AddedProducts map[string]*struct {
	ProductID string `json:"productID"`
	Count     int    `json:"count"`
}

type CartUsecase interface {
	AddProductToCart(ctx context.Context, u *User, productID string, count int) error
	RemoveProductFromCart(ctx context.Context, u *User, cartKey string) error
	ShowUsersCart(ctx context.Context, u *User) (*Cart, error)
}

type CartRepository interface {
	AddToCart(ctx context.Context, u *User, productID string, count int) error
	RemoveFromCart(ctx context.Context, u *User, cartKey string) error
	GetUsersCart(ctx context.Context, u *User) (*Cart, error)
}

func NewCart() *Cart {
	c := &Cart{
		AddedProducts: AddedProducts{},
	}
	return c
}
