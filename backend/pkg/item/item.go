package item

type Item struct {
	ID string `json:"id"` // id товара, вычисляется из характеристики, может быть sku для внешнийх систем

	Type string `json:"type"` // вид товара, стартер генератор или еще что то
	Name string `json:"name"` // oem имя товара
	SKU  string `json:"sku"`  // артикул товара
	Char string `json:"char"` // характеристика товара

	Images []Image `json:"images"`

	Brandcars string `json:"brandcars"` // Марка
	Modelcars string `json:"modelcars"` // Модель
	Engine    string `json:"engine"`    // двигатель
	Year      string `json:"year"`      // год

	Amount        int `json:"amount"`
	Price         int `json:"price"`         // цена без учета скидок наценок
	DiscountPrice int `json:"discountPrice"` // цена со скидками наценками
	Discount      int `json:"discount"`      // процент или сумма скидки, как понять цена или процент?
}

type Image struct {
	Owner string `json:"owner"`
	Path  string `json:"path"` // путь без хоста
	Main  bool   `json:"main"` // основное фото
}
