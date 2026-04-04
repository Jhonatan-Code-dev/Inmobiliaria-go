package storage

import (
	"context"
	"io"
	"time"
)

type Storage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string) error
	GetURL(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
