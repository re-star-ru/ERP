package main

import (
	"backend/cmd/proxy/item"
	"backend/pkg/item/delivery"
	"backend/pkg/item/repo"
	"backend/pkg/item/usecase"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func main() {
	log.Logger = log.
		With().Caller().
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		fs := http.FileServer(http.Dir("/home/restar/git/erp/site/public"))
		r.Handle("/*", fs)
	})

	r.Get("/item", item.Serve)

	// - 1c авторизация
	// s3 put get delete
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	s := &S3{minioClient, "srv1c"}
	r.Route("/s3", func(s3r chi.Router) {
		s3r.Use(middleware.Logger)

		ctx := context.Background()

		err = minioClient.MakeBucket(ctx, s.srv1cbucket, minio.MakeBucketOptions{})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, s.srv1cbucket)
			if errBucketExists == nil && exists {
				log.Printf("we already own %s", s.srv1cbucket)
			} else {
				log.Fatal().Err(err).Send()
			}
		} else {
			log.Printf("suscessfully created %s", s.srv1cbucket)
		}

		s3r.Put("/image", s.PutImage)
		s3r.Delete("/image", s.DeleteImage)
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

		r.Use(middleware.Logger)
		r.Get("/products", func(w http.ResponseWriter, r *http.Request) {
			ps, err := c.Products(100)
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

	log.Debug().Msg("listen at :" + os.Getenv("HOST"))
	log.Fatal().Err(http.ListenAndServe(":"+os.Getenv("HOST"), r))
}

func logError(w http.ResponseWriter, err error, statusCode int, errorInfo string) {
	log.Err(err).Msg(errorInfo)
	http.Error(w, err.Error(), statusCode)
}

type S3 struct {
	client      *minio.Client
	srv1cbucket string
}

func (s *S3) PutImage(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join("images", uuid.NewV4().String()+".jpeg")
	log.Debug().Msg(filePath)
	info, err := s.client.PutObject(
		context.Background(), s.srv1cbucket, filePath, r.Body, -1,
		minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")},
	)

	if err != nil {
		log.Err(err).Msg("ошибка загрузки изображения")
	}

	log.Printf("путь к файлу: %s/%s", info.Bucket, info.Key)

	w.Write([]byte(path.Join(info.Bucket, info.Key)))
}

func (s *S3) DeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("пока не реализованно"))
}
