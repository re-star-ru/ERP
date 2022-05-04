package img

import (
	"backend/pkg"
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

// PutImage put image to s3 and return path bucket/key
func (s *ImageService) PutImage(w http.ResponseWriter, r *http.Request) {

	// str, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("cant read body")
	// 	return
	// }
	// log.Error().Msgf("body: %s", str)

	if r.ContentLength < 200 {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("not an image"), "wrong content length")
		return
	}

	filePath := path.Join("images", uuid.NewString()+".jpeg")

	info, err := s.m.PutObject(
		context.Background(), s.srv1cbucket, filePath, r.Body, r.ContentLength,
		minio.PutObjectOptions{ContentType: r.Header.Get("Content-Type")},
	)
	if err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error image upload")
		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, path.Join(info.Bucket, info.Key))
}

func (s *ImageService) DeleteImage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		pkg.SendErrorJSON(w, r, http.StatusBadRequest, errors.New("id is empty"), "id is empty")
		return
	}

	if err := s.m.RemoveObject(r.Context(), s.srv1cbucket, "images/"+id, minio.RemoveObjectOptions{}); err != nil {
		pkg.SendErrorJSON(w, r, http.StatusInternalServerError, err, "error image delete")
		return
	}

	render.Status(r, http.StatusOK)
}
