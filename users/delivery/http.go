package delivery

import "net/http"

type Repo interface {
}

type HttpDelivery struct {
}

func NewHttpDelivery(r Repo) *HttpDelivery {
	return &HttpDelivery{}
}

// creating new user
func (h *HttpDelivery) Register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (h *HttpDelivery) Rbac(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (h *HttpDelivery) GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}
