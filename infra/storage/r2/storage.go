package r2

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	client     *s3.Client
	bucketName string
	accountID  string
	publicURL  string
}

func NewStorageService(client *s3.Client, bucketName, accountID, publicURL string) *StorageService {
	return &StorageService{
		client:     client,
		bucketName: bucketName,
		accountID:  accountID,
		publicURL:  publicURL,
	}
}

func (s *StorageService) Upload(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("error al subir objeto a R2: %w", err)
	}
	return nil
}

func (s *StorageService) GetURL(ctx context.Context, key string) (string, error) {
	// Usar la URL pública personalizada (ej: https://archivos.cuotamax.com/imagenes/foto.jpg)
	return fmt.Sprintf("%s/%s", s.publicURL, key), nil
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error al eliminar objeto de R2: %w", err)
	}
	return nil
}

func (s *StorageService) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignResult, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", fmt.Errorf("error al generar presigned URL: %w", err)
	}

	return presignResult.URL, nil
}
