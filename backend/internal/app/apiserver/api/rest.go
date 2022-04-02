package api

import (
	"backend/internal/app/apiserver/api/auth"
	"backend/internal/app/apiserver/api/catalog"
	"backend/internal/app/apiserver/api/users"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/docgen"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi"
)

// Rest is rest server struct
type Rest struct {
	Version       string
	Authenticator *auth.Service

	httpServer *http.Server
	lock       sync.Mutex
}

// Run is run rest server
func (s *Rest) Run(port int) {
	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(port, s.routes())
	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()
	log.Println("http server terminated", err)
}

// Shutdown is shutdown server
func (s *Rest) Shutdown() {
	log.Println("shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Println("http shutdown error", err)
		}
		log.Println("shutdown http server completed")
	}

	s.lock.Unlock()

}

func (s *Rest) routes() chi.Router {
	//imagesint
	imgInit()
	r := chi.NewRouter()

	var corsMiddleware = cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(corsMiddleware.Handler)

	// auth
	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(5 * time.Second))
		r.Use(middleware.SetHeader("Content-type", "application/json"))

		r.Mount("/auth/login", s.Authenticator.LoginHandler())
		r.Mount("/auth/registration", s.Authenticator.RegistrationHandler())
	})

	// api v1, def protected, use casbin for auth
	r.Route("/v1", func(rapi chi.Router) {
		//middlewares of route
		rapi.Use(middleware.SetHeader("Content-type", "application/json"))
		rapi.Use(middleware.Timeout(60 * time.Second))
		rapi.Use(s.Authenticator.Authenticator)
		rapi.Use(s.Authenticator.Authorizer)

		// TODO: delete user by email
		// user
		rapi.Route("/user", func(r chi.Router) {
			r.Get("/", users.GetUser)                // get user info
			r.Put("/", users.Update)                 // update user info
			r.Put("/password", users.UpdatePassword) // update user password
		})
		rapi.Get("/users", users.GetAllUsers)

		//offers
		rapi.Get("/offers", getOffers)
		rapi.Get("/offer/{GUID}", getOfferByGUID)

		//catalog
		rapi.Route("/catalog", func(rproducts chi.Router) {
			rproducts.Get("/update", catalog.UpdateProductStore)
			rproducts.Get("/info", catalog.GetInfo)
			rproducts.Get("/product-types", catalog.GetProductTypes)

			rproducts.Get("/*", catalog.GetDefaultGroupList)
			//rproducts.Get("/page-{page}", catalog.GetDefaultGroupList)

			//rproducts.Get("/list/by-sku", catalog.GetAllSKUList)
			//rproducts.Get("/list/by-sku/*", catalog.GetProductListParams)
			//rproducts.Get("/list/by-sku/{productTypeGuid}/{offset}", catalog.GetSKUListByTypeAndOffset)
		})

		//images
		rapi.Post("/image/{GUID}", uploadImage)
		rapi.Delete("/image/{GUID}", deleteImage)
	})

	r.With(middleware.SetHeader("Content-type", "application/json")).Get("/health-check", healthCheckHandler)
	r.With(middleware.SetHeader("Content-type", "application/json")).Get("/doc", func(w http.ResponseWriter, req *http.Request) {
		if _, err := w.Write([]byte(docgen.JSONRoutesDoc(r))); err != nil {
			log.Println(err)
		}
	})

	// swagger api
	r.Group(func(sr chi.Router) {
		//sr.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger"))))
		sr.Handle("/swagger/*", http.StripPrefix("/swagger/", http.FileServer(http.Dir("swagger"))))

	})

	//docgen.PrintRoutes(router)
	return r
}

func (s *Rest) makeHTTPServer(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, `{"alive": true}`); err != nil {
		log.Println(err)
	}
}
