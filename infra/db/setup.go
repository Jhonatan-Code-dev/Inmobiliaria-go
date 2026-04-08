package db

import (
	"rentals-go/ent"
)

func Setup(dsn string) (*ent.Client, error) {
	database := NewDB(dsn)

	client, err := database.GetClient()
	if err != nil {
		return nil, err
	}

	if err := database.Migrate(); err != nil {
		return nil, err
	}

	// Datos iniciales básicos
	if err := seedRoles(client); err != nil {
		return nil, err
	}
	if err := seedAdmin(client); err != nil {
		return nil, err
	}
	if err := SeedTiposPago(client); err != nil {
		return nil, err
	}

	return client, nil
}
