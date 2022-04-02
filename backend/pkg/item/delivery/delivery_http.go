package delivery

import (
	"backend/pkg"
	"backend/pkg/item"
	"encoding/json"
	"github.com/go-chi/chi"
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
	if len(p) < 3 {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, nil, "необходимо минимум 3 символа для поиска") // todo
		return
	}

	ps, err := it.iu.Search(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ps)

	log.Debug().Msgf("limit: %s", r.URL.Query().Get("limit"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if err := it.iu.UploadPricelists(limit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Err(err).Msg("cant update price list")
		return
	}
}
