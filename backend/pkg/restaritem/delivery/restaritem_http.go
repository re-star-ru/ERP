package delivery

import (
	"backend/pkg/restaritem"
	"github.com/go-chi/render"
	"net/http"
)

// будет json ручка
// адрес для заполнения данных по товару
// адрес для выкладывания товара на сайт
// список всех товаров на сайте так же

type IHTTPRestaritemUsecase interface {
	Create(restaritem *restaritem.RestarItem) (*restaritem.RestarItem, error)
	GetAll() ([]restaritem.RestarItem, error) // pagination?
	GetByID(id int) (*restaritem.RestarItem, error)
}

func NewHTTPRestaritemDelivery(uc IHTTPRestaritemUsecase) *HTTPRestaritemDelivery {
	return &HTTPRestaritemDelivery{
		uc: uc,
	}
}

type HTTPRestaritemDelivery struct {
	uc IHTTPRestaritemUsecase
}

// 1: создать новый итем, возвращает id итема, который потом надо перенаправить в qr код
func (h *HTTPRestaritemDelivery) Create(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"message": "not implemented"})
}

// 2: получить все итемы
func (h *HTTPRestaritemDelivery) GetAll(w http.ResponseWriter, r *http.Request) {

}

// 3: получить итем по id
func (h *HTTPRestaritemDelivery) GetByID(w http.ResponseWriter, r *http.Request) {

}

// 4: страница с данными об итеме, включая дефекты, работы, фото
func (h *HTTPRestaritemDelivery) RestaritemView(w http.ResponseWriter, r *http.Request) {

}
