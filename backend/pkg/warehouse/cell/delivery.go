package cell

import (
	"backend/pkg"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"html/template"
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

	cellProducts, err := c.us.ByID(cellID)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, err, "cant get cell")

		return
	}

	if r.URL.Query().Has("json") {
		render.JSON(w, r, pkg.JSON{"body": cellProducts})

		return
	}

	files := []string{
		"./web/template/cell.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant parse template")

		return
	}

	if err = tmpl.Execute(w, cellProducts); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "cant execute template")

		return
	}
}
