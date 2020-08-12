package authDeliveryHTTP

import (
	"backend/internal/app/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{us}
	authEndpoints := e.Group("auth")
	{
		authEndpoints.POST("/sign-up", handler.SignUp)
		authEndpoints.POST("/sign-in", handler.SignIn)
	}
}

func (u *UserHandler) SignUp(c echo.Context) (err error) {
	var user domain.User
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
	var user domain.User
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

func isRequestValid(m *domain.User) (bool, error) {
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
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrUserAlreadyExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
