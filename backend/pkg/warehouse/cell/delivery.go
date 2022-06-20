package cell

import (
	"backend/pkg"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

type ICellUsecase interface {
	ByID(cellID string) (*Cell, error)
}

type Delivery struct {
	us ICellUsecase
}

func NewCellDelivery(uc ICellUsecase) *Delivery {
	return &Delivery{uc}
}

var ErrCellNotFound = errors.New("cell not found")

func (c *Delivery) Get(w http.ResponseWriter, r *http.Request) {
	cellID := chi.URLParam(r, "cellID")
	if cellID == "" {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, ErrCellNotFound, "cellID is required")

		return
	}

	cells, err := c.us.ByID(cellID)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get cell")

		return
	}

	render.JSON(w, r, pkg.JSON{"body": cells})
}
