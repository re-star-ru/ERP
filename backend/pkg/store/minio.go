package store

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type MinioStore struct {
	minio        *minio.Client
	endpoint     string
	bucket       string
	bucketPolicy string
}

func NewMinioStore(c *minio.Client) (stor *MinioStore, err error) {
	stor = &MinioStore{minio: c, bucket: "oprox",
		endpoint: "https://s3.re-star.ru",
		bucketPolicy: `
	{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": [
						"*"
					]
				},
				"Action": [
					"s3:GetBucketLocation",
					"s3:ListBucket",
					"s3:ListBucketMultipartUploads"
				],
				"Resource": [
					"arn:aws:s3:::oprox"
				]
			},
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": [
						"*"
					]
				},
				"Action": [
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::oprox/*"
				]
			}
		]
	}`}

	err = stor.minio.MakeBucket(context.Background(), stor.bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := stor.minio.BucketExists(context.Background(), stor.bucket)
		if errBucketExists == nil && exists {
			log.Info().Msgf("minio already own %s", stor.bucket)
		} else {
			return nil, err
		}
	} else {
		log.Info().Msgf("minio bucket suscessfully created %s", stor.bucket)
	}

	if err = stor.minio.SetBucketPolicy(context.Background(), stor.bucket, stor.bucketPolicy); err != nil {
		return nil, err
	}

	return stor, nil
}

func (m *MinioStore) Store(fpath, contentType string, r io.Reader) (string, error) {
	info, err := m.minio.PutObject(context.Background(), m.bucket, fpath, r, -1,
		minio.PutObjectOptions{ContentType: contentType},
	)

	return path.Join(info.Bucket, info.Key), err
}

func (m *MinioStore) Path() string {
	return fmt.Sprintf("%s/%s", m.endpoint, m.bucket)
}
