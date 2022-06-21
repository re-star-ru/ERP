package delivery

// будет json ручка
// адрес для заполнения данных по товару
// адрес для выкладывания товара на сайт
// список всех товаров на сайте так же

type IHTTPRestaritemUsecase interface {
}

func NewHTTPRestaritemDelivery(uc IHTTPRestaritemUsecase) *HTTPRestaritemDelivery {
	return &HTTPRestaritemDelivery{
		uc: uc,
	}
}

type HTTPRestaritemDelivery struct {
	uc IHTTPRestaritemUsecase
}
