package productDeliveryHTTP

import (
	"backend/internal/app/domain"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

type productHandler struct {
	usecase domain.ProductUsecase
}

func NewHandler(e *echo.Group, pu domain.ProductUsecase) {
	handler := &productHandler{pu}

	e.GET("", handler.Get)
	e.POST("", handler.Create)
}

type getResponse struct {
	Products []*domain.Product `json:"products"`
}

func (p *productHandler) Get(c echo.Context) error {
	ps, err := p.usecase.GetProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, &getResponse{ps})
}

type createInput struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

func toProduct(c createInput, u *domain.User) *domain.Product {
	return &domain.Product{
		Name:         c.Name,
		SKU:          c.SKU,
		Creator:      *u,
		CreatedAt:    time.Now(),
		LastModified: time.Now(),
	}
}

func (p *productHandler) Create(c echo.Context) error {
	inp := &createInput{}

	if err := c.Bind(inp); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}

	user := c.Get(domain.UserKey).(*domain.User)
	logrus.Println("getting user from ctx:", user)

	pr := toProduct(*inp, user)

	if err := p.usecase.CreateProduct(c.Request().Context(), user, pr); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}

	return p.respond(c, http.StatusCreated, pr)
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
