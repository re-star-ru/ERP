package productDeliveryHTTP

import (
	"backend/internal/app/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: update one product
func (p *productHandler) update(c echo.Context) error {
	upd := models.Product{}

	if err := c.Bind(&upd); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}
	user := p.getUserFromCtx(c)

	if err := p.usecase.UpdateProduct(c.Request().Context(), user, upd); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}
	return p.respond(c, http.StatusNoContent, nil)
}
