package db

import (
	"context"
	"log"

	"rentals-go/config/security"
	"rentals-go/ent"
	"rentals-go/ent/admin"
	"rentals-go/ent/rol"
)

// seedAdmin crea el super admin inicial si no existe.
func seedAdmin(client *ent.Client) error {
	ctx := context.Background()
	const (
		superUsuario = "Jhonatan"
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
		SetHashContrasena(hash).
		Save(ctx); err != nil {
		return err
	}
	log.Println("🟦 Super admin creado por defecto")
	return nil
}

// seedRoles inicializa los roles obligatorios del sistema.
func seedRoles(client *ent.Client) error {
	ctx := context.Background()
	const (
		rolAdmin       = "administrador"
		rolAdminDesc   = "Rol con acceso administrativo a la empresa"
	)

	exists, err := client.Rol.
		Query().
		Where(rol.NombreEQ(rolAdmin)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if _, err := client.Rol.
		Create().
		SetNombre(rolAdmin).
		SetDescripcion(rolAdminDesc).
		Save(ctx); err != nil {
		return err
	}
	log.Println("🟦 Rol 'administrador' inicializado")
	return nil
}
