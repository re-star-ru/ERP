package productDeliveryHTTP

import (
	"backend/internal/app/models"
	"backend/internal/app/product"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type productHandler struct {
	usecase product.Usecase
}

func NewHandler(e *echo.Group, pu product.Usecase) {
	handler := &productHandler{pu}

	e.GET("", handler.Get)
	e.POST("", handler.Create)
	e.PUT("/:id", handler.update)
	e.GET("/:id", handler.GetOneByID)
}

type getResponse struct {
	Products []models.Product `json:"products"`
}

type getProduct struct {
	models.Product
}

func (p *productHandler) Get(c echo.Context) error {
	ps, err := p.usecase.GetProducts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, &getResponse{ps})
}

type inputGetByID struct {
	ID string `json:"id"`
}

func (p *productHandler) GetOneByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return p.error(c, http.StatusUnprocessableEntity, models.ErrBadParamInput)
	}

	pd, err := p.usecase.GetOne(c.Request().Context(), id)
	if err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, &getProduct{pd})
}

type createInput struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

func toProduct(c createInput) models.Product {
	return models.Product{
		Name: c.Name,
		SKU:  c.SKU,
	}
}

func (p *productHandler) Create(c echo.Context) error {
	inp := createInput{}

	if err := c.Bind(&inp); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}

	user := p.getUserFromCtx(c)

	pr := toProduct(inp)

	if err := p.usecase.CreateProduct(c.Request().Context(), user, pr); err != nil {
		return p.error(c, http.StatusBadRequest, err)
	}

	return p.respond(c, http.StatusCreated, "ok")
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

func (p *productHandler) getUserFromCtx(c echo.Context) *models.User {
	user := c.Get(models.UserKey).(*models.User)
	logrus.Println("getting user from ctx:", user)
	return user
}
