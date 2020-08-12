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
		if authHeader == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			return c.NoContent(http.StatusUnauthorized)
		}

		if headerParts[0] != "Bearer" {
			return c.NoContent(http.StatusUnauthorized)
		}

		user, err := m.usecase.ParseToken(c.Request().Context(), headerParts[1])
		if err != nil {
			status := http.StatusInternalServerError
			if err == domain.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}

			return c.NoContent(status)
		}

		logrus.Println("setting user in ctx:", user)
		c.Set(domain.UserKey, user)

		return next(c)

		//email, passwordHash, ok := r.BasicAuth()

		// пропускаем анонимных пользователей дальше
		//if !ok {
		//	log.Println("пустые данные аутентификации, будет предоставлен анонимный доступ")
		//	ctx := context.WithValue(r.Context(), "aclGroup", "anonymous")
		//	ctx = context.WithValue(ctx, "email", "anonymous")
		//	next.ServeHTTP(w, r.WithContext(ctx))
		//	return next(c)
		//}
		//
		//u, err := users.GetUserByEmail(email)
		//if err != nil {
		//	err := errors.New("не существует такого пользователя, перенаправление на аутентификацию")
		//	log.Println(err)
		//	log.Println(email, passwordHash)
		//	http.Error(w, err.Error(), http.StatusUnauthorized)
		//	return next(c)
		//}
		//
		//// перенаправляем на авторизацию если хеш не совпадает с тем что в базе
		//if !u.CheckUsersPasswordHash(passwordHash) {
		//	err := errors.New("неверные данные аутентификации")
		//	log.Println(err)
		//	log.Println(email, passwordHash)
		//	http.Error(w, err.Error(), http.StatusUnauthorized)
		//	return next(c)
		//}
		//
		//ctx := context.WithValue(r.Context(), "aclGroup", u.AclGroup)
		//ctx = context.WithValue(ctx, "email", u.Email)
		//
	}
}

func InitMiddleware(userUC domain.UserUsecase) *GoMiddleware {
	return &GoMiddleware{
		usecase: userUC,
	}
}
