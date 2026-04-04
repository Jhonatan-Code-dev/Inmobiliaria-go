package valkey

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func New(
	addr string,
	defaultTTL time.Duration,
) (*Repository, func(), error) {

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Valkey no disponible (%s): %v", addr, err)
	} else {
		log.Printf("✅ Valkey conectado (%s)", addr)
	}

	repo := &Repository{
		client:     rdb,
		defaultTTL: defaultTTL,
	}

	cleanup := func() {
		_ = rdb.Close()
	}

	return repo, cleanup, nil
}
