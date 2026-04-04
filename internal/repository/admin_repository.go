package repository

import (
	"context"

	"rentals-go/ent"
	"rentals-go/ent/admin"
	"rentals-go/internal/domain"
)

type AdminRepoEnt struct {
	client *ent.Client
}

func NewAdminRepo(client *ent.Client) *AdminRepoEnt {
	return &AdminRepoEnt{client: client}
}

func (r *AdminRepoEnt) BuscarPorUsuario(ctx context.Context, usuario string) (*domain.Admin, error) {
	a, err := r.client.Admin.Query().
		Where(admin.UsuarioEQ(usuario)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Admin{
		ID:             a.ID,
		Nombre:         a.Nombre,
		Usuario:        a.Usuario,
		HashContrasena: a.HashContrasena,
		Activo:         a.Activo,
	}, nil
}

func (r *AdminRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Admin, error) {
	a, err := r.client.Admin.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Admin{
		ID:             a.ID,
		Nombre:         a.Nombre,
		Usuario:        a.Usuario,
		HashContrasena: a.HashContrasena,
		Activo:         a.Activo,
	}, nil
}

func (r *AdminRepoEnt) ActualizarCredenciales(ctx context.Context, id int, usuarioNuevo, hashContrasena string) (*domain.Admin, error) {
	a, err := r.client.Admin.UpdateOneID(id).
		SetUsuario(usuarioNuevo).
		SetHashContrasena(hashContrasena).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Admin{
		ID:             a.ID,
		Nombre:         a.Nombre,
		Usuario:        a.Usuario,
		HashContrasena: a.HashContrasena,
		Activo:         a.Activo,
	}, nil
}
