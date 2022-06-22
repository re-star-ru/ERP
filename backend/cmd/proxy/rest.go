package main

import (
	"backend/configs"
	"backend/ent"
	"backend/pkg/oneclient"
	"backend/pkg/qr"
	restaritemDelivery "backend/pkg/restaritem/delivery"
	restaritemRepo "backend/pkg/restaritem/repo"
	restaritemUsecase "backend/pkg/restaritem/usecase"
	"backend/pkg/warehouse/cell"
	"context"
	"errors"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"

	"backend/pkg"
	"backend/pkg/img"
	"backend/pkg/item/delivery"
	"backend/pkg/item/repo"
	"backend/pkg/item/usecase"
	"backend/pkg/pricelist"
	"backend/pkg/store"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/lib/pq"
)

//nolint:funlen
func Rest(config configs.Config) *chi.Mux {
	log.Debug().
		Str("MINIO", config.Endpoint).
		Str("ONEC", config.OnecHost).
		Msg("resources endpoints")

	minioClient := newMinio(config)

	stor, err := store.NewMinioStore(minioClient) // global app bucket
	if err != nil {
		log.Fatal().Err(err).Msgf("cant create minio store")

		return nil
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	if !config.Production {
		log.Debug().Msg("development...")
		router.Use(cors.AllowAll().Handler)
	}

	// TODO: Authorized routes and anonymouse route
	router.Route("/s3", func(s3r chi.Router) {
		is := img.NewImageService(minioClient, "srv1c") // srv1c image bucket
		s3r.Put("/image", is.PutImage)
		// s3r.With(SimpleAuthMiddleware).Put("/image", is.PutImage)
		s3r.With(SimpleAuthMiddleware).Delete("/image/{id}", is.DeleteImage)
	})
	// -

	onecClient := oneclient.NewOneClient(config.OnecHost, config.OnecToken)

	itemRepo := repo.NewRepoOnec(onecClient)
	itemUsecase := usecase.NewItemUsecase(itemRepo)
	itemHTTP := delivery.NewItemDelivery(itemUsecase)

	pricer := pricelist.NewPricerUsecase(stor, itemRepo)
	priceSrv := pricelist.NewPricelistHttp(pricer)

	{
		// site api
		router.Get("/search/{query}", itemHTTP.SearchBySKU)
		router.Get("/catalog", itemHTTP.CatalogHandler)
	}

	{
		// restar item - инфа по товарам в ремонте или приемке
		client := initEnt(config.PG)

		rirepo := restaritemRepo.NewRestaritemRepo(client)
		riUcase := restaritemUsecase.NewRestaritemUsecase(rirepo)
		riDelivery := restaritemDelivery.NewHTTPRestaritemDelivery(riUcase)

		router.Post("/restaritem", riDelivery.Create)
		router.Get("/restaritem", riDelivery.GetAll)
		router.Get("/restaritem/{id}", riDelivery.RestaritemView)
	}

	{
		// qr code сервис, получает параметры -
		// текстовую строку для печати qr кода
		// x, y - ширина высота картинки,
		// и какой то шаблон для печати либо один из заранее подготовленных шаблонов
		// я пока не решил
		// возвращает готовый к печати pdf файл в формате а4,
		// в левом верхнем углну находится сам qr с параметрами

		qrd := qr.NewHTTPDelivery()
		router.Post("/qr", qrd.NewQRCode)

		cellRepo := cell.NewRepoOnec(onecClient)
		cellUC := cell.NewCellUsecase(cellRepo)
		cellDel := cell.NewCellDelivery(cellUC)

		router.Get("/warehouse/cell/{cellID}", cellDel.Get)
	}

	{
		// pricelist api
		router.Route("/pricelists", func(r chi.Router) {
			r.Get("/", priceSrv.PricelistHandler)
			r.Post("/refresh", priceSrv.ManualRefreshHandler)
			r.Post("/meily", priceSrv.MeiliRequest)
			r.Get("/{name}", priceSrv.PricelistByConsumerHandler)
		})
	}

	return router
}

var ErrNotAuthorized = errors.New("not authorized")

func SimpleAuthMiddleware(next http.Handler) http.Handler {
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal().Msg("API_KEY env not set")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("SimpleAuthMiddleware")

		if r.Header.Get("X-API-KEY") != apiKey {
			pkg.SendErrorJSON(w, r, http.StatusUnauthorized,
				ErrNotAuthorized, "wrong api key: "+r.Header.Get("X-API-KEY"),
			)

			return
		}

		log.Debug().Msg("auth ok")
		next.ServeHTTP(w, r)
	})
}

func newMinio(config configs.Config) *minio.Client {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Minio init error")
	}

	return minioClient
}

func initEnt(path string) *ent.Client {
	client, err := ent.Open("postgres", path)
	if err != nil {
		log.Fatal().Err(err).Msg("ent init error")
	}

	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("failed to create schema")
	}

	return client
}
