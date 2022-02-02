package main

import (
	"backend/internal/app/apiserver/s3"
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
		s := &S3{}
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

		err = minioClient

		bucketName := "srv1c"

		s3r.Put("/image", s.PutImage)
	})
	// -
	http.ListenAndServe(":8000", r)
}

type S3 struct {
}

func (s *S3) PutImage(r http.ResponseWriter, w *http.Request) {
	s3.New()

	s3.UploadFile()

}
