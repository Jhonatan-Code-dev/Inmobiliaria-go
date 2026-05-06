package db

import (
	"context"
	"log"

	"rentals-go/config/security"
	"rentals-go/ent"
	"rentals-go/ent/admin"
	"rentals-go/ent/rol"
	"rentals-go/ent/tipoidentificacion"
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

	roles := []struct {
		nombre      string
		descripcion string
	}{
		{nombre: "administrador", descripcion: "Rol con acceso administrativo total a la empresa"},
		{nombre: "supervisor", descripcion: "Rol para supervisión operativa y reportes"},
		{nombre: "vendedor", descripcion: "Rol para gestión comercial, clientes y contratos"},
		{nombre: "inventario", descripcion: "Rol para control de activos y estado de unidades"},
	}

	for _, r := range roles {
		exists, err := client.Rol.
			Query().
			Where(rol.NombreEQ(r.nombre)).
			Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		if _, err := client.Rol.
			Create().
			SetNombre(r.nombre).
			SetDescripcion(r.descripcion).
			Save(ctx); err != nil {
			return err
		}
		log.Printf("🟦 Rol '%s' inicializado\n", r.nombre)
	}
	return nil
}

func seedTiposIdentificacion(client *ent.Client) error {
	ctx := context.Background()

	type item struct {
		codigo string
		nombre string
		pais   *string
	}

	pe := "PE"
	items := []item{
		{codigo: "DNI", nombre: "Documento Nacional de Identidad", pais: &pe},
		{codigo: "CE", nombre: "Carnet de Extranjeria", pais: &pe},
		{codigo: "PAS", nombre: "Pasaporte", pais: nil},
		{codigo: "RUC", nombre: "Registro Unico de Contribuyentes", pais: &pe},
	}

	for _, it := range items {
		exists, err := client.TipoIdentificacion.
			Query().
			Where(tipoidentificacion.CodigoEQ(it.codigo)).
			Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		builder := client.TipoIdentificacion.
			Create().
			SetCodigo(it.codigo).
			SetNombre(it.nombre).
			SetActivo(true)
		if it.pais != nil {
			builder.SetPais(*it.pais)
		}

		if _, err := builder.Save(ctx); err != nil {
			return err
		}
	}

	log.Println("🟦 Tipos de identificación inicializados")
	return nil
}
