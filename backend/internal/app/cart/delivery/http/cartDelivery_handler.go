package cartDeliveryHTTP

import (
	"backend/internal/app/domain"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

type cartHandler struct {
	usecase domain.CartUsecase
}

func NewHandler(e *echo.Group, cu domain.CartUsecase) {
	handler := &cartHandler{cu}

	e.GET("", handler.Get)

}

func (h cartHandler) Get(c echo.Context) error {
	user := c.Get(domain.UserKey).(*domain.User)

	logrus.Println("getting user from ctx:", user)

	cart, err := h.usecase.ShowUsersCart(c.Request().Context(), user)
	if err != nil {
		return h.error(c, http.StatusBadRequest, err)
	}

	return h.respond(c, http.StatusOK, cart)
}

func (h *cartHandler) error(c echo.Context, code int, err error) error {
	return h.respond(c, code, map[string]string{"error": err.Error()})
}

func (h *cartHandler) respond(c echo.Context, code int, data interface{}) error {
	if data != nil {
		if err := c.JSON(code, data); err != nil {
			logrus.Errorln(err)
			return err
		}
	}
	if err := c.NoContent(code); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
