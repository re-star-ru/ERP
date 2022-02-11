package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	r := chi.NewRouter()
	r.With()

	// - 1c авторизация
	// s3 put get delete
	r.Route("/s3", func(s3r chi.Router) {
		s3r.Use(middleware.Logger)

		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKey := os.Getenv("MINIO_ACCESS_KEY")
		secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
			Secure: false,
		})
		if err != nil {
			log.Fatalln(err)
		}

		s := &S3{minioClient, "srv1c"}

		ctx := context.Background()

		err = minioClient.MakeBucket(ctx, s.srv1cbucket, minio.MakeBucketOptions{})
		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, s.srv1cbucket)
			if errBucketExists == nil && exists {
				log.Printf("we already own %s\n", s.srv1cbucket)
			} else {
				log.Fatalln(err)
			}
		} else {
			log.Printf("suscessfully created %s\n\n", s.srv1cbucket)
		}

		s3r.Put("/image", s.PutImage)
		s3r.Delete("/image", s.DeleteImage)
	})
	// -
	log.Fatalln(http.ListenAndServe(":8000", r))
}

type S3 struct {
	client      *minio.Client
	srv1cbucket string
}

func (s *S3) PutImage(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join("images", uuid.NewV4().String()+".jpeg")
	log.Println(filePath)
	info, err := s.client.PutObject(
		context.Background(), s.srv1cbucket, filePath, r.Body, -1,
		minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")},
	)

	if err != nil {
		log.Println("ошибка загрузки изображения: ", err)
	}

	log.Printf("путь к файлу: %s/%s \n", info.Bucket, info.Key)

	w.Write([]byte(path.Join(info.Bucket, info.Key)))
}

func (s *S3) DeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("пока не реализованно"))
}
