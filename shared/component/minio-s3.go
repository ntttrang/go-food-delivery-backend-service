package shareComponent

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3Uploader struct {
	apiKey     string
	bucketName string
	domain     string
	region     string
	secretKey  string
	client     *minio.Client
	useSSL     bool
}

func NewS3Uploader(apiKey, bucketName, domain, region, secretKey string, useSSL bool) (*s3Uploader, error) {
	minioClient, err := minio.New(domain, &minio.Options{
		Creds:  credentials.NewStaticV4(apiKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &s3Uploader{
		apiKey:     apiKey,
		bucketName: bucketName,
		domain:     domain,
		region:     region,
		secretKey:  secretKey,
		client:     minioClient,
	}, nil
}

func (s *s3Uploader) SaveFileUpload(ctx context.Context, filename string, filePath string, contentType string) error {

	info, err := s.client.FPutObject(ctx, s.bucketName, filename, filePath,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		log.Printf("upload err: %v", err)
		return err
	}
	log.Printf("info: %v", info)
	return nil
}

func (c *s3Uploader) GetDomain() string {
	return c.domain
}
