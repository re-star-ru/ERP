package productDeliveryHTTP

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: update one product
func (p *productHandler) update(c echo.Context) error {
	ps, err := p.usecase.GetProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, &getResponse{ps})
}
