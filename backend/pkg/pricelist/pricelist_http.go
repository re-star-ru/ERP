package pricelist

import (
	"backend/pkg"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type Pricer interface {
	GetPricelistByConsumerName(string) (string, error) // return path to s3 pricelist by consumer name
	GetPricelists() (map[string]string, error)         // return map [consumer: pricelist]
	Update() error                                     // request update pricelists
}

type PricelistHttpService struct {
	pricer Pricer
}

func NewPricelistHttp(p Pricer) *PricelistHttpService {
	return &PricelistHttpService{
		pricer: p,
	}
}

func (ps *PricelistHttpService) ManualRefreshHandler(w http.ResponseWriter, r *http.Request) {

	if err := ps.pricer.Update(); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant update")
		return
	}

}

func (ps *PricelistHttpService) MeiliRequest(w http.ResponseWriter, r *http.Request) {

	log.Println("meili request")

}

func (ps *PricelistHttpService) PricelistHandler(w http.ResponseWriter, r *http.Request) {
	prices, err := ps.pricer.GetPricelists()
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cannot get pricelists")
		return
	}

	render.JSON(w, r, prices)
}

func (it *PricelistHttpService) PricelistByConsumerHandler(w http.ResponseWriter, r *http.Request) {

}
