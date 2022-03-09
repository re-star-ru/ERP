package delivery

import "net/http"

type ItemDelivery struct {
}

func NewItemDelivery() *ItemDelivery {
	return &ItemDelivery{}
}

func (it *ItemDelivery) UpdatePricelists(w http.ResponseWriter, r *http.Request) {

	

}
