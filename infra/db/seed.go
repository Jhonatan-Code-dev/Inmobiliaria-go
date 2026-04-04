package db

import (
	"context"
	"log"

	"rentals-go/config/security"
	"rentals-go/ent"
	"rentals-go/ent/admin"
)

// seedAdmin crea el super admin inicial si no existe.
func seedAdmin(client *ent.Client) error {
	ctx := context.Background()
	const (
		superUsuario = "Jhonatan"
		superNombre  = "Super Admin"
		superPass    = "912059555"
	)

	exists, err := client.Admin.
		Query().
		Where(admin.UsuarioEQ(superUsuario)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	hasher := security.NewServicioHash()
	hash, err := hasher.Encriptar(superPass)
	if err != nil {
		return err
	}

	if _, err := client.Admin.
		Create().
		SetUsuario(superUsuario).
		SetNombre(superNombre).
		SetHashContrasena(hash).
		Save(ctx); err != nil {
		return err
	}
	log.Println("🟦 Super admin creado por defecto")
	return nil
}
