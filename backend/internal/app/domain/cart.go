package domain

import "context"

type Cart struct {
	OwnerID       string `json:"userID"`
	AddedProducts map[string]*struct {
		ProductID string `json:"productID"`
		Count     int    `json:"count"`
	} `json:"addedProducts"`
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
