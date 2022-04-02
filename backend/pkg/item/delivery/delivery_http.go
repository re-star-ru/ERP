package delivery

import (
	"backend/pkg"
	"backend/pkg/item"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type IItemUsecase interface {
	UploadPricelists(limit int) error

	Search(s string) (map[string]item.Item, error)
}

type ItemDelivery struct {
	iu IItemUsecase
}

func NewItemDelivery(iu IItemUsecase) *ItemDelivery {
	return &ItemDelivery{
		iu: iu,
	}
}

func (it *ItemDelivery) UpdatePricelists(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msgf("limit: %s", r.URL.Query().Get("limit"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if err := it.iu.UploadPricelists(limit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Err(err).Msg("cant update price list")
		return
	}
}

func (it *ItemDelivery) SearchBySKU(w http.ResponseWriter, r *http.Request) {
	p := chi.URLParam(r, "query")
	log.Printf("chi %v", p)
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
