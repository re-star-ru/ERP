package delivery

import (
	"backend/pkg"
	"backend/pkg/restaritem"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"html/template"
	"net/http"
	"strconv"
)

// будет json ручка
// адрес для заполнения данных по товару
// адрес для выкладывания товара на сайт
// список всех товаров на сайте так же

type IHTTPRestaritemUsecase interface {
	Create(ctx context.Context, restaritem restaritem.RestarItem) (*restaritem.RestarItem, error)
	GetAll(ctx context.Context) ([]*restaritem.RestarItem, error) // pagination?
	GetByID(ctx context.Context, id int) (*restaritem.RestarItem, error)

	AddPhoto(ctx context.Context, id int, photo []byte) error
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
	ritem := &restaritem.RestarItem{}
	if err := json.NewDecoder(r.Body).Decode(ritem); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse json")

		return
	}

	nitem, err := h.uc.Create(r.Context(), *ritem)
	if err != nil {
		pkg.SendErrorJSON(w, r, pkg.StatuscodeByError(err), err, "cant create restaritem")

		return
	}

	render.JSON(w, r, nitem)
}

// 2: получить все итемы
func (h *HTTPRestaritemDelivery) GetAll(w http.ResponseWriter, r *http.Request) {
	ritems, err := h.uc.GetAll(r.Context())
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant get restaritems")

		return
	}

	render.JSON(w, r, ritems)
}

// 4: страница с данными об итеме, включая дефекты, работы, фото
func (h *HTTPRestaritemDelivery) RestaritemView(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "cant parse id")

		return
	}

	ritem, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get restaritem")

		return
	}

	if r.URL.Query().Has("json") {
		render.JSON(w, r, ritem)

		return
	}

	tmpl, err := template.ParseFiles("./web/template/restaritem.html")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	if err = tmpl.Execute(w, ritem); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}

func parseId(r *http.Request) (int, error) {
	sid := chi.URLParam(r, "id")
	if sid == "" {
		return 0, pkg.ErrWrongInput
	}

	id, err := strconv.Atoi(sid)
	if err != nil {
		return 0, pkg.ErrWrongInput
	}

	return id, nil
}

func (h *HTTPRestaritemDelivery) AddPhoto(w http.ResponseWriter, r *http.Request) {
	id, err := parseId(r)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse id")

		return
	}

	ritem, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get restaritem")

		return
	}

	// load image
	// render 5 sizes of image
	// save images to s3,
	// return Image{} struct
	// update ritem with this Image{} struct
	// return new ritem

	render.JSON(w, r, ritem)
}
