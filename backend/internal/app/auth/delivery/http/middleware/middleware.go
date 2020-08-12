package authDeliveryHTTPMiddleware

import (
	"backend/internal/app/domain"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo"
)

type GoMiddleware struct {
	usecase domain.UserUsecase
}

func (m *GoMiddleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		return next(c)
	}
}

// Authenticator уверяется что пользователь тот за кого себя выдает и устанавливает пустого пользователя если
// нет учетных данных
func (m *GoMiddleware) Authenticator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		// If auth header empty set default anonymous user to context
		if authHeader == "" {
			logrus.Println("setting empty user in ctx:", domain.User{})
			c.Set(domain.UserKey, domain.User{})
			return next(c)
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return m.error(c, http.StatusUnauthorized, domain.ErrUnauthorized)
		}

		if headerParts[0] != "Bearer" {
			return m.error(c, http.StatusUnauthorized, domain.ErrUnauthorized)
		}

		user, err := m.usecase.ParseToken(c.Request().Context(), headerParts[1])
		if err != nil {
			status := http.StatusInternalServerError
			if err == domain.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}
			return m.error(c, status, err)
		}

		logrus.Println("setting user in ctx:", user)
		c.Set(domain.UserKey, user)

		return next(c)
	}
}

func InitMiddleware(userUC domain.UserUsecase) *GoMiddleware {
	return &GoMiddleware{
		usecase: userUC,
	}
}

func (m *GoMiddleware) error(c echo.Context, code int, err error) error {
	return m.respond(c, code, map[string]string{"error": err.Error()})
}

func (m *GoMiddleware) respond(c echo.Context, code int, data interface{}) error {
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
