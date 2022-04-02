package models

type Cart struct {
	OwnerID       string `json:"userID"`
	AddedProducts `json:"addedProducts"`
}

type AddedProducts map[string]*struct {
	ProductID string `json:"productID"`
	Count     int    `json:"count"`
}

func NewCart() *Cart {
	c := &Cart{
		AddedProducts: AddedProducts{},
	}
	return c
}
