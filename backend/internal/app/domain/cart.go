package domain

import "context"

type Cart struct {
	UserID        string `json:"userID"`
	AddedProducts map[string]struct {
		ProductID string `json:"productID"`
		Count     int    `json:"count"`
	} `json:"addedProducts"`
}

type CartUsecase interface {
	AddToCart(ctx context.Context, u *User, productID string, count int) error
	RemoveFromCart(ctx context.Context, u *User, cartKey string) error
}

type CartRepository interface {
	AddToCart(ctx context.Context, u *User, productID string, count int) error
	RemoveFromCart(ctx context.Context, u *User, cartKey string) error
}
