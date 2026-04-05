package repository

import (
	"context"

	"rentals-go/ent"
	"rentals-go/ent/usuario"
	"rentals-go/internal/domain"
)

type UsuarioRepoEnt struct {
	client *ent.Client
}

func NewUsuarioRepo(client *ent.Client) *UsuarioRepoEnt {
	return &UsuarioRepoEnt{client: client}
}

func (r *UsuarioRepoEnt) Crear(ctx context.Context, u *domain.Usuario) (*domain.Usuario, error) {
	usr, err := r.client.Usuario.Create().
		SetUsuario(u.Usuario).
		SetHashContrasena(u.HashContrasena).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Usuario{
		ID:             usr.ID,
		Usuario:        usr.Usuario,
		HashContrasena: usr.HashContrasena,
		Estado:         usr.Estado,
	}, nil
}

func (r *UsuarioRepoEnt) BuscarPorUsuario(ctx context.Context, username string) (*domain.Usuario, error) {
	u, err := r.client.Usuario.
		Query().
		Where(usuario.UsuarioEQ(username)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Usuario{
		ID:             u.ID,
		Usuario:        u.Usuario,
		HashContrasena: u.HashContrasena,
		Estado:         u.Estado,
	}, nil
}

func (r *UsuarioRepoEnt) BuscarPerfil(ctx context.Context, id int) (*domain.Usuario, *domain.Empresa, error) {
	u, err := r.client.Usuario.
		Query().
		Where(usuario.IDEQ(id)).
		WithEmpresasUsuario(func(q *ent.EmpresaUsuarioQuery) {
			q.WithEmpresa()
		}).
		First(ctx)
	if err != nil {
		return nil, nil, err
	}
	emp := &domain.Empresa{}
	if len(u.Edges.EmpresasUsuario) > 0 && u.Edges.EmpresasUsuario[0].Edges.Empresa != nil {
		e := u.Edges.EmpresasUsuario[0].Edges.Empresa
		emp = &domain.Empresa{
			ID:              e.ID,
			Nombre:          e.Nombre,
			Pais:            ptrToString(e.Pais),
			Moneda:          e.Moneda,
		}
	}
	return &domain.Usuario{
		ID:             u.ID,
		Usuario:        u.Usuario,
		HashContrasena: u.HashContrasena,
		Estado:         u.Estado,
	}, emp, nil
}
