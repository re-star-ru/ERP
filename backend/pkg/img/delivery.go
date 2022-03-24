package img

import (
	"context"
	"net/http"
	"path"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

func NewImageService(m *minio.Client, srv1ImageBucket string) *ImageService {
	ctx := context.Background()

	// TODO: create with standart politics
	err := m.MakeBucket(ctx, srv1ImageBucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := m.BucketExists(ctx, srv1ImageBucket)
		if errBucketExists == nil && exists {
			log.Printf("we already own %s", srv1ImageBucket)
		} else {
			log.Fatal().Err(err).Send()
		}
	} else {
		log.Printf("suscessfully created %s", srv1ImageBucket)
	}

	return &ImageService{
		m:           m,
		srv1cbucket: srv1ImageBucket,
	}
}

type ImageService struct {
	m           *minio.Client
	srv1cbucket string
}

func (s *ImageService) PutImage(w http.ResponseWriter, r *http.Request) {
	filePath := path.Join("images", uuid.NewString()+".jpeg")
	log.Debug().Msg(filePath)
	info, err := s.m.PutObject(
		context.Background(), s.srv1cbucket, filePath, r.Body, -1,
		minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")},
	)

	if err != nil {
		log.Err(err).Msg("ошибка загрузки изображения")
	}

	log.Printf("путь к файлу: %s/%s", info.Bucket, info.Key)

	w.Write([]byte(path.Join(info.Bucket, info.Key)))
}

func (s *ImageService) DeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("пока не реализованно"))
}
