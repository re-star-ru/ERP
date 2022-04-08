package delivery

import (
	"backend/pkg"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IItemUsecase interface {
	Search(s string) ([]interface{}, error)
}

type ItemDelivery struct {
	iu IItemUsecase
}

func NewItemDelivery(iu IItemUsecase) *ItemDelivery {
	return &ItemDelivery{
		iu: iu,
	}
}

func (it *ItemDelivery) SearchBySKU(w http.ResponseWriter, r *http.Request) {
	p := chi.URLParam(r, "query")
	if len(p) < 3 {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "необходимо минимум 3 символа для поиска")
		return
	}

	ps, err := it.iu.Search(p)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "ошибка поиска")
		return
	}

	if err = json.NewEncoder(w).Encode(ps); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "ошибка кодировки")
		return
	}
}

func (it *ItemDelivery) CatalogHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if len(q) < 3 {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, pkg.ErrWrongInput, "необходимо минимум 3 символа для поиска")
		return
	}

	ps, err := it.iu.Search(q)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "ошибка поиска")
		return
	}

	if err = json.NewEncoder(w).Encode(ps); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "ошибка кодировки")
		return
	}
}
