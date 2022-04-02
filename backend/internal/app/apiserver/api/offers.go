package api

import (
	"backend/internal/app/apiserver/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func getOffers(w http.ResponseWriter, _ *http.Request) {
	offers, err := db.GetOffersFrom1C()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rData, err := json.Marshal(offers)
	if err != nil {
		log.Fatalln(err)
	}
	if code, err := w.Write(rData); err != nil {
		log.Println(code, err)
	}
}

func getOfferByGUID(w http.ResponseWriter, r *http.Request) {
	guid := chi.URLParam(r, "GUID")
	offer, err := db.GetOffer(guid)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	renderImgPaths(&offer)
	rData, err := json.Marshal(offer)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(rData); err != nil {
		log.Println(err)
	}
}

func renderImgPaths(o *db.Offer) {
	if o.MainImage != "" {
		o.MainImage = fmt.Sprintf("%v%v", db.Offers.S3path, o.MainImage)
	}
	if len(o.Images) == 0 {
		o.Images = []string{}
		return
	}
	for i := range o.Images {
		o.Images[i] = fmt.Sprintf("%v%v", db.Offers.S3path, o.Images[i])
	}
}
