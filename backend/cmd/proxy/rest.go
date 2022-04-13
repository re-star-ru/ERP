package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"backend/pkg/img"
	"backend/pkg/item/delivery"
	"backend/pkg/item/repo"
	"backend/pkg/item/usecase"
	"backend/pkg/pricelist"
	"backend/pkg/store"
)

type cfg struct {
	endpoint, accessKey, secretAccessKey string
	onecHost, onecToken                  string
}

func newMinio(c cfg) *minio.Client {
	minioClient, err := minio.New(c.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.accessKey, c.secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Minio init error")
	}

	return minioClient
}

func Rest(c cfg) *chi.Mux {
	log.Info().Str("MINIO", c.endpoint).Str("ONEC", c.onecHost).Msg("resourse endpoints")

	minioClient := newMinio(c)
	stor, err := store.NewMinioStore(minioClient) // global app bucket
	if err != nil {
		log.Fatal().Err(err).Msgf("cant create minio store")
		return nil
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler) // todo: if dev

	// TODO: Authorized routes and anonymouse route
	r.Route("/s3", func(s3r chi.Router) {
		is := img.NewImageService(minioClient, "srv1c") // srv1c image bucket
		s3r.Use(middleware.Logger)
		s3r.Put("/image", is.PutImage)
		s3r.Delete("/image", is.DeleteImage)
	})
	// -

	itemRepo := repo.NewRepoOnec(c.onecHost, c.onecToken)
	itemUsecase := usecase.NewItemUsecase(itemRepo)
	itemHttp := delivery.NewItemDelivery(itemUsecase)

	pricer := pricelist.NewPricerUsecase(stor, itemRepo)
	priceSrv := pricelist.NewPricelistHttp(pricer)

	{
		// site api
		r.Get("/search/{query}", itemHttp.SearchBySKU)
		r.Get("/catalog", itemHttp.CatalogHandler)
	}

	{
		// pricelist api
		r.Route("/pricelists", func(r chi.Router) {
			r.Get("/", priceSrv.PricelistHandler)
			r.Get("/{name}", priceSrv.PricelistByConsumerHandler)
			r.Post("/refresh", priceSrv.ManualRefreshHandler)
		})
	}

	return r
}
