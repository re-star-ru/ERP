package s3

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gabriel-vasile/mimetype"

	"github.com/gofrs/uuid"

	"github.com/minio/minio-go/v6"
)

var s3ClientConfig = struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	location        string
	useSSL          bool
}{
	"s3.restar26.site",
	"cYeEFXYTeX",
	"cYeEFXYTeX",
	"us-east-1",
	true,
}

type s3ClientType struct {
	Client *minio.Client
}

var s3Client = s3ClientType{
	Client: initClient(),
}

const bucketName = "public"

func (s s3ClientType) MakeBucket(bucketName string) {
	err := s.Client.MakeBucket(bucketName, s3ClientConfig.location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := s.Client.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func initClient() *minio.Client {
	c, err := minio.New(s3ClientConfig.endpoint, s3ClientConfig.accessKeyID, s3ClientConfig.secretAccessKey, s3ClientConfig.useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func UploadFile(file []byte, mime *mimetype.MIME) (path string) {
	u, err := uuid.NewV4()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("sizebefore:  ", int64(len(file)))
	log.Println(u.String() + mime.Extension())
	size, err := s3Client.Client.PutObject(bucketName, u.String()+mime.Extension(), bytes.NewReader(file), int64(len(file)), minio.PutObjectOptions{
		ContentType: mime.String(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("sizeafter ", size)
	return fmt.Sprintf("/%v/%v%v", bucketName, u.String(), mime.Extension())
}

func S3Init() {

	// Initialize minio client object.

	//log.Printf("%#v\n", s3Client.Client) // minioClient is now setup

	objectName := "test/testsss"
	filePath := "test.png"
	contentType := "image/png"

	n, err := s3Client.Client.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	//objectName2 := "test2.png"
	filePath2 := "test2.png"

	if err := s3Client.Client.FGetObject(bucketName, objectName, filePath2, minio.GetObjectOptions{}); err != nil {
		log.Println("Not downloaded by", err)
	}
}
