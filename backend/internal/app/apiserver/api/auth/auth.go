package auth

import (
	"backend/internal/app/apiserver/api/users"
	"context"
	"errors"
	"log"
	"net/http"

	. "github.com/casbin/casbin"
)

var e = NewEnforcer("configs/acl/acl_model.conf", "configs/acl/acl_policy.csv")

type Service struct {
}

func (s *Service) LoginHandler() http.Handler {
	return http.HandlerFunc(users.Login)
}

func (s *Service) RegistrationHandler() http.Handler {
	return http.HandlerFunc(users.Register)
}

// Authenticator уверяется что пользователь тот за кого себя выдает и устанавливает пустого пользователя если
// нет учетных данных
func (s *Service) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth middleware")
		email, passwordHash, ok := r.BasicAuth()

		// пропускаем анонимных пользователей дальше
		if !ok {
			log.Println("пустые данные аутентификации, будет предоставлен анонимный доступ")
			ctx := context.WithValue(r.Context(), "aclGroup", "anonymous")
			ctx = context.WithValue(ctx, "email", "anonymous")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		u, err := users.GetUserByEmail(email)
		if err != nil {
			err := errors.New("не существует такого пользователя, перенаправление на аутентификацию")
			log.Println(err)
			log.Println(email, passwordHash)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// перенаправляем на авторизацию если хеш не совпадает с тем что в базе
		if !u.CheckUsersPasswordHash(passwordHash) {
			err := errors.New("неверные данные аутентификации")
			log.Println(err)
			log.Println(email, passwordHash)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "aclGroup", u.AclGroup)
		ctx = context.WithValue(ctx, "email", u.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/// Authorizer проверяет есть ли у пользователя доступ к методу по пути
func (s *Service) Authorizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aclGroup := r.Context().Value("aclGroup").(string)
		method := r.Method
		path := r.URL.Path
		log.Printf(`sub "%v" obj "%v" act "%v"`, aclGroup, path, method)
		log.Printf("enforce is %v", e.Enforce(aclGroup, path, method))
		if e.Enforce(aclGroup, path, method) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, http.StatusText(403), http.StatusForbidden)
		}
	})
}
