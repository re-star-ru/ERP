package cartDeliveryHTTP

import (
	"backend/internal/app/domain"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

type cartHandler struct {
	usecase domain.CartUsecase
}

func NewHandler(e echo.Group, cu domain.CartUsecase) {
	handler := &cartHandler{cu}

	e.GET("", handler.Get)
}

func (h cartHandler) Get(c echo.Context) error {
	user := c.Get(domain.UserKey).(*domain.User)

	logrus.Println("getting user from ctx:", user)

	cart, err := h.usecase.ShowUsersCart(c.Request().Context(), user)
	if err != nil {
		return err
	}

	return
}

func (p *productHandler) error(c echo.Context, code int, err error) error {
	return p.respond(c, code, map[string]string{"error": err.Error()})
}

func (p *productHandler) respond(c echo.Context, code int, data interface{}) error {
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
