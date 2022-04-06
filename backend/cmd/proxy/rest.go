package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/cors"
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
	log.Info().Str("MINIO", c.endpoint).Str("ONEC", c.onecHost).Send()

	minioClient, err := minio.New(c.endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(c.accessKey, c.secretAccessKey, ""),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Minio init error")
	}

	is := img.NewImageService(minioClient, "srv1c")

	r := chi.NewRouter()
	//r.Use(cors.Default().Handler)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"*"},
	}).Handler)

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

	itemRepo := repo.NewRepoOnec(c.onecHost, c.onecToken)
	itemUsecase := usecase.NewItemUsecase(itemRepo, minioClient)
	itemHttp := delivery.NewItemDelivery(itemUsecase)

	r.Get("/search/{query}", itemHttp.SearchBySKU)
	r.Get("/catalog", itemHttp.CatalogHandler)

	r.Route("/1c", func(r chi.Router) {
		r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
			ps, err := itemRepo.Items(100, 100)
			if err != nil {
				logError(w, err, 400, "cant get products")
				return
			}

			if err := json.NewEncoder(w).Encode(ps); err != nil {
				logError(w, err, 500, "cannot unmarshal")
				return
			}

		})

		r.Post("/updatePricelist", itemHttp.UpdatePricelists)

	})

	return r
}
