package authDeliveryHTTP

import (
	"backend/internal/app/auth"
	"backend/internal/app/models"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserUsecase auth.Usecase
}

func NewHandler(authEndpoints *echo.Group, us auth.Usecase) {
	handler := &UserHandler{us}
	{
		authEndpoints.POST("/sign-up", handler.SignUp)
		authEndpoints.POST("/sign-in", handler.SignIn)
		authEndpoints.GET("/whoami", handler.whoami)

	}
}

func (u *UserHandler) SignUp(c echo.Context) (err error) {
	var user models.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok && err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	err = u.UserUsecase.SignUp(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) SignIn(c echo.Context) (err error) {
	var user models.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok && err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()

	var token string
	token, err = u.UserUsecase.SignIn(ctx, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{err.Error()})
	}

	return c.JSON(http.StatusOK, token)
}

func (u *UserHandler) whoami(c echo.Context) (err error) {
	user := c.Get(models.UserKey).(*models.User)

	return c.JSON(http.StatusOK, user)
}

func isRequestValid(m *models.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
