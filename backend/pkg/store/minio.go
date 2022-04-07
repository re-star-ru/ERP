package store

import (
	"context"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
)

type MinioStore struct {
	minio        *minio.Client
	bucket       string
	bucketPolicy string
}

func NewMinioStore(c *minio.Client) (*MinioStore, error) {
	m := &MinioStore{minio: c, bucket: "oprox", bucketPolicy: `
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

	err := m.minio.MakeBucket(context.Background(), m.bucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := m.minio.BucketExists(context.Background(), m.bucket)
		if errBucketExists == nil && exists {
			log.Info().Msgf("minio already own %s", m.bucket)
		} else {
			return nil, err
		}
	} else {
		log.Info().Msgf("minio bucket suscessfully created %s", m.bucket)
	}

	if err := m.minio.SetBucketPolicy(context.Background(), m.bucket, m.bucketPolicy); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *MinioStore) Store(fpath, contentType string, r io.Reader) (string, error) {
	info, err := m.minio.PutObject(context.Background(), m.bucket, fpath, r, -1,
		minio.PutObjectOptions{ContentType: contentType},
	)

	return path.Join(info.Bucket, info.Key), err
}
