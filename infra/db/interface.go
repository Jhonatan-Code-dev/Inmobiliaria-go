package db

import (
	"context"

	"rentals-go/ent"
)

type Database interface {
	GetClient() (*ent.Client, error)
	Ping(ctx context.Context) error
	Migrate() error
	Close() error
}
