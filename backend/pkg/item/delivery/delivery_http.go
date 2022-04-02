package delivery

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type IItemUsecase interface {
	UploadPricelists(limit int) error
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
	log.Debug().Msgf("limit: %s", r.URL.Query().Get("limit"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if err := it.iu.UploadPricelists(limit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Err(err).Msg("cant update price list")
		return
	}
}
