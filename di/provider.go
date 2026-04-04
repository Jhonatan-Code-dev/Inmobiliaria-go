package di

import (
	"rentals-go/config/env"
	"rentals-go/ent"
	"rentals-go/infra/db"

	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(ProvideConfig)
var DatabaseSet = wire.NewSet(ProvideEntClient)

func ProvideConfig() (*env.Config, error) {
	cfg := env.NewConfig()
	return cfg, nil
}

func ProvideEntClient(cfg *env.Config) (*ent.Client, error) {
	return db.Setup(cfg.DSN)
}

func ProvideJWTSecret(cfg *env.Config) []byte {
	return []byte(cfg.JWTSecret)
}
