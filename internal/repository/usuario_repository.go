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
		SetNombres(u.Nombres).
		SetNillableApellidos(nilIfEmpty(u.Apellidos)).
		SetCorreo(u.Correo).
		SetNillableTelefono(nilIfEmpty(u.Telefono)).
		SetHashContrasena(u.HashContrasena).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Usuario{
		ID:             usr.ID,
		Nombres:        usr.Nombres,
		Apellidos:      ptrToString(usr.Apellidos),
		Correo:         usr.Correo,
		Telefono:       ptrToString(usr.Telefono),
		HashContrasena: usr.HashContrasena,
	}, nil
}

func (r *UsuarioRepoEnt) BuscarPorCorreo(ctx context.Context, correo string) (*domain.Usuario, error) {
	u, err := r.client.Usuario.
		Query().
		Where(usuario.CorreoEQ(correo)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.Usuario{
		ID:             u.ID,
		Nombres:        u.Nombres,
		Apellidos:      ptrToString(u.Apellidos),
		Correo:         u.Correo,
		Telefono:       ptrToString(u.Telefono),
		HashContrasena: u.HashContrasena,
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
			DocumentoFiscal: ptrToString(e.DocumentoFiscal),
			Correo:          ptrToString(e.Correo),
			Telefono:        ptrToString(e.Telefono),
			Direccion:       ptrToString(e.Direccion),
			Ciudad:          ptrToString(e.Ciudad),
			Pais:            ptrToString(e.Pais),
			Moneda:          e.Moneda,
		}
	}
	return &domain.Usuario{
		ID:             u.ID,
		Nombres:        u.Nombres,
		Apellidos:      ptrToString(u.Apellidos),
		Correo:         u.Correo,
		Telefono:       ptrToString(u.Telefono),
		HashContrasena: u.HashContrasena,
	}, emp, nil
}
