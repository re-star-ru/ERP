package cell

import "net/http"

type CellDelivery struct {
}

func NewCellDelivery() *CellDelivery {
	return &CellDelivery{}
}

func (c *CellDelivery) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get cell"))
}
