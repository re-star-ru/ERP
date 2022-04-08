package item

type Item struct {
	ID   string `json:"id"`
	SKU  string `json:"sku"`
	Name string `json:"name"`
	Char string `json:"char"`

	Images []struct {
		Owner string `json:"owner"`
		Path  string `json:"path"`
		Main  bool   `json:"main"`
	} `json:"images"`

	Amount        int `json:"amount"`
	Price         int `json:"price"`
	DiscountPrice int `json:"discountPrice"`
	Discount      int `json:"discount"`
}
