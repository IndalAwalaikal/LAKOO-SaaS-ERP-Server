package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService interface {
	UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error)
	GetFileUrl(objectName string) string
}

type minioService struct {
	client *minio.Client
	bucket string
}

func NewMinioService(endpoint, accessKey, secretKey, bucket string, useSSL bool) MinioService {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Make sure bucket exists
	ctx := context.Background()
	exists, errBucketExists := minioClient.BucketExists(ctx, bucket)
	if errBucketExists == nil && !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Failed to create bucket %s: %v", bucket, err)
		} else {
			log.Printf("Successfully created bucket %s", bucket)
			
			// Optional: Set bucket policy to public read so media files can be loaded via URL directly
			policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": ["arn:aws:s3:::` + bucket + `/*"],"Sid": ""}]}`
			err = minioClient.SetBucketPolicy(ctx, bucket, policy)
			if err != nil {
				log.Printf("Warning: Failed to set public bucket policy: %v", err)
			}
		}
	} else if errBucketExists != nil {
		log.Printf("Warning: Checking bucket existence failed: %v", errBucketExists)
	}

	return &minioService{
		client: minioClient,
		bucket: bucket,
	}
}

func (s *minioService) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	// Return the relative minio path (which can be resolved by client)
	return "/" + s.bucket + "/" + objectName, nil
}

func (s *minioService) GetFileUrl(objectName string) string {
	return "/" + s.bucket + "/" + objectName
}
