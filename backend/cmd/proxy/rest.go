package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"backend/cmd/proxy/item"
	"backend/pkg/img"
	"backend/pkg/item/delivery"
	"backend/pkg/item/repo"
	"backend/pkg/item/usecase"
)

type cfg struct {
	endpoint, accessKey, secretAccessKey string
	onecHost, onecToken                  string
}

func Rest(c cfg) *chi.Mux {
	minioClient, err := minio.New(c.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.accessKey, c.secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	is := img.NewImageService(minioClient, "srv1c")

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		fs := http.FileServer(http.Dir("/home/restar/git/erp/site/public")) // wtf?
		r.Handle("/*", fs)
	})

	r.Get("/item", item.Serve)
	r.Route("/s3", func(s3r chi.Router) {
		s3r.Use(middleware.Logger)
		s3r.Put("/image", is.PutImage)
		s3r.Delete("/image", is.DeleteImage)
	})

	// -

	r.Route("/1c", func(r chi.Router) {
		c := repo.NewClient1c(
			os.Getenv("ONEC_HOST"),
			os.Getenv("ONEC_TOKEN"),
		)
		id := delivery.NewItemDelivery(
			usecase.NewItemUsecase(c, minioClient),
		)

		r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
			ps, err := c.Products(100, 100)
			if err != nil {
				logError(w, err, 400, "cant get products")
				return
			}

			if err := json.NewEncoder(w).Encode(ps); err != nil {
				logError(w, err, 500, "cannot unmarshal")
				return
			}

		})
		r.Post("/updatePricelist", id.UpdatePricelists)

	})

	return r
}
