package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()

	// - 1c авторизация
	// s3 put get delete
	r.Route("/s3", func(s3r chi.Router) {
		endpoint := "node2.re-star.ru"
		accessKey := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
			Secure: false,
		})
		if err != nil {
			log.Fatalln(err)
		}

		s := &S3{minioClient}

		ctx := context.Background()

		bucketName := "srv1c"

		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		s3r.Put("/image", s.PutImage)
	})
	// -
	http.ListenAndServe(":8000", r)
}

type S3 struct {
	client *minio.Client
}

func (s *S3) PutImage(r http.ResponseWriter, w *http.Request) {

}
