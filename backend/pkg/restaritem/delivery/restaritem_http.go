package delivery

import (
	"backend/pkg"
	"backend/pkg/photo"
	"backend/pkg/restaritem"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

// будет json ручка
// адрес для заполнения данных по товару
// адрес для выкладывания товара на сайт
// список всех товаров на сайте так же

type IRestaritemUsecase interface {
	Create(ctx context.Context, restaritem restaritem.RestarItem) (*restaritem.RestarItem, error)
	GetAll(ctx context.Context) ([]*restaritem.RestarItem, error) // pagination?
	GetByID(ctx context.Context, id int) (*restaritem.RestarItem, error)

	AddPhoto(ctx context.Context, id int, photo photo.Photo) (*restaritem.RestarItem, error)
}

type IPhotoUsecase interface {
	NewPhoto(ctx context.Context, dir string, photo io.ReadCloser) (photo.Photo, error)
}

func NewHTTPRestaritemDelivery(uc IRestaritemUsecase, phuc IPhotoUsecase) *HTTPRestaritemDelivery {
	return &HTTPRestaritemDelivery{
		uc:   uc,
		phuc: phuc,
	}
}

type HTTPRestaritemDelivery struct {
	uc   IRestaritemUsecase
	phuc IPhotoUsecase
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
	id, err := parseInt(r, "id")
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

	files := []string{
		"./web/template/restaritem.html",
		"./web/template/photos_view.html",
		"./web/template/inspection_view.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	if err = tmpl.Execute(w, ritem); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}

func parseInt(r *http.Request, key string) (int, error) {
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
	id, err := parseInt(r, "id")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse id")

		return
	}

	_, err = h.uc.GetByID(r.Context(), id)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get restaritem")

		return
	}

	//if err = r.ParseMultipartForm(); err != nil {
	//	pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse form")
	//
	//	return
	//}

	// load image
	f, ff, err := r.FormFile("photo")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get file")

		return
	}

	log.Printf("input photo Content-Type: %+v", ff.Header.Get("Content-Type"))

	nPhoto, err := h.phuc.NewPhoto(r.Context(), "restaritem", f)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant create photo")

		return
	}

	nritem, err := h.uc.AddPhoto(r.Context(), id, nPhoto)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant add photo")

		return
	}

	log.Printf("ritem: %+v", nritem)

	// load image
	// render 5 sizes of image
	// save images to s3,
	// return Image{} struct
	// update ritem with this Image{} struct
	// return new ritem

	files := []string{
		"./web/template/photos_view.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	if err = tmpl.Execute(w, nritem); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}

func (h *HTTPRestaritemDelivery) SetInspectionByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(r, "id")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse id")

		return
	}

	inspID := chi.URLParam(r, "inspectiondID")
	if inspID == "" {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "cant parse inspectiondID")

		return
	}

	rating, err := parseInt(r, "rating")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse rating")

		return
	}

	_, err = h.uc.GetByID(r.Context(), id)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get restaritem")

		return
	}

	files := []string{
		"./web/template/inspection_view.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	a, ok := starterInspections[inspID]
	if !ok {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "cant find inspection")

		return
	}

	a.Quality = rating
	starterInspections[inspID] = a

	if err = tmpl.Execute(w, pkg.JSON{"id": id, "inspections": starterInspections}); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}

func (h *HTTPRestaritemDelivery) ListInspections(w http.ResponseWriter, r *http.Request) {
	id, err := parseInt(r, "id")
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant parse id")

		return
	}

	_, err = h.uc.GetByID(r.Context(), id)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get restaritem")

		return
	}

	files := []string{
		"./web/template/inspection_view.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	if err = tmpl.Execute(w, pkg.JSON{"id": id, "inspections": starterInspections}); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}

var starterInspections = map[string]struct {
	Name    string
	Quality int
}{
	"1": {
		"Бендикс",
		-1, // не определено по умолчанию
	},
	"2": {
		"Вилка стартера",
		-1,
	},
	"3": {
		"Втулка стартера",
		-1,
	},
	"4": {
		"Втягивающее реле",
		-1,
	},
	"5": {
		"Дополнительное реле",
		-1,
	},
	"6": {
		"Крышка стартера",
		-1,
	},
	"7": {
		"Муфта стартера",
		-1,
	},
	"8": {
		"Планетарный механизм",
		-1,
	},
	"9": {
		"Подшипник стартера",
		-1,
	},
	"11": {
		"Статор стартера",
		-1,
	},
	"12": {
		"Щётки стартера",
		-1,
	},
	"13": {
		"Щеточный узел стартера",
		-1,
	},
	"14": {
		"Якорь",
		-1,
	},
}

var generatorInspections = []string{
	"Вакуумный насос",
	"Диодный мост",
	"Диоды",
	"Коллектор",
	"Крышка генератора",
	"Подшипник генератора",
	"Проставка под подшипник",
	"Реле регулятор",
	"Ротор",
	"Сальник",
	"Статор генератора",
	"Шкив",
	"Щётки генератора",
	"Щёточный узел генератора",
}
