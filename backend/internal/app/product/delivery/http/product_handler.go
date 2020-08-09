package productDeliveryHTTP

import (
	"backend/internal/app/domain"
	"net/http"

	"github.com/labstack/echo"
)

type productHandler struct {
	usecase domain.ProductUsecase
}

func NewHandler(e *echo.Echo, pu domain.ProductUsecase) {
	handler := &productHandler{pu}

	e.GET("/products", handler.Get)

	e.POST("/products", handler.Create)
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
}

func toProduct(c createInput) *domain.Product {
	return &domain.Product{
		Name: c.Name,
	}
}

func (p *productHandler) Create(c echo.Context) error {
	inp := new(createInput)

	if err := c.Bind(inp); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	//user := c.MustGet(auth.CtxUserKey).(*auth.User)
	pr := toProduct(*inp)

	if err := p.usecase.CreateProduct(c.Request().Context(), &domain.User{}, pr); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, pr)
}
