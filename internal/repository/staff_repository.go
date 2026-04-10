package repository

import (
	"context"
	"fmt"
	"rentals-go/ent"
	"rentals-go/ent/empresausuario"
	"rentals-go/ent/usuario"
	"rentals-go/internal/domain"
)

type StaffRepoEnt struct {
	client *ent.Client
}

func NewStaffRepo(client *ent.Client) *StaffRepoEnt {
	return &StaffRepoEnt{client: client}
}

func (r *StaffRepoEnt) Listar(ctx context.Context, filtros domain.StaffFiltros) ([]*domain.Staff, int, error) {
	query := r.client.EmpresaUsuario.Query().
		Where(empresausuario.EmpresaID(filtros.EmpresaID)).
		WithUsuario().
		WithRol()

	if filtros.Busqueda != "" {
		query = query.Where(empresausuario.HasUsuarioWith(usuario.UsuarioContains(filtros.Busqueda)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	entities, err := query.Limit(filtros.PorPagina).Offset(offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.Staff
	for _, e := range entities {
		results = append(results, mapStaffToDomain(e))
	}

	return results, total, nil
}

func (r *StaffRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Staff, error) {
	e, err := r.client.EmpresaUsuario.Query().
		Where(empresausuario.IDEQ(id), empresausuario.EmpresaIDEQ(empresaID)).
		WithUsuario().
		WithRol().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapStaffToDomain(e), nil
}

func (r *StaffRepoEnt) Crear(ctx context.Context, s *domain.RegistroStaff, hash string) (*domain.Staff, error) {
	// Transacción para crear Usuario y vincularlo a Empresa
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	u, err := tx.Usuario.Create().
		SetUsuario(s.Usuario).
		SetHashContrasena(hash).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	eu, err := tx.EmpresaUsuario.Create().
		SetEmpresaID(s.EmpresaID).
		SetUsuarioID(u.ID).
		SetRolID(s.RolID).
		SetEstado(empresausuario.EstadoActivo).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Recargar con edges
	return r.BuscarPorID(ctx, eu.ID, s.EmpresaID)
}

func (r *StaffRepoEnt) Actualizar(ctx context.Context, id int, empresaID int, rolID int, estado string) (*domain.Staff, error) {
	err := r.client.EmpresaUsuario.Update().
		Where(empresausuario.IDEQ(id), empresausuario.EmpresaIDEQ(empresaID)).
		SetRolID(rolID).
		SetEstado(empresausuario.Estado(estado)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return r.BuscarPorID(ctx, id, empresaID)
}

func (r *StaffRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	eu, err := r.client.EmpresaUsuario.Query().
		Where(empresausuario.IDEQ(id), empresausuario.EmpresaIDEQ(empresaID)).
		Only(ctx)
	if err != nil {
		return err
	}
	
	if eu.Principal {
		return fmt.Errorf("no se puede eliminar al usuario principal")
	}

	return r.client.EmpresaUsuario.DeleteOneID(id).Exec(ctx)
}

func mapStaffToDomain(e *ent.EmpresaUsuario) *domain.Staff {
	s := &domain.Staff{
		ID:        e.ID,
		UsuarioID: e.UsuarioID,
		RolID:     e.RolID,
		EmpresaID: e.EmpresaID,
		Principal: e.Principal,
		Estado:    string(e.Estado),
	}
	if e.Edges.Usuario != nil {
		s.Usuario = e.Edges.Usuario.Usuario
	}
	if e.Edges.Rol != nil {
		s.RolNombre = e.Edges.Rol.Nombre
	}
	return s
}
