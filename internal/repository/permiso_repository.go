package repository

import (
	"context"
	"fmt"

	"rentals-go/ent"
	"rentals-go/ent/permiso"
	"rentals-go/internal/domain"
)

type PermisoRepoEnt struct {
	client *ent.Client
}

func NewPermisoRepo(client *ent.Client) *PermisoRepoEnt {
	return &PermisoRepoEnt{client: client}
}

func (r *PermisoRepoEnt) ListarPaginado(ctx context.Context, filtros domain.PermisoFiltros) ([]*domain.Permiso, int, error) {
	query := r.client.Permiso.Query().
		Where(permiso.EmpresaIDEQ(filtros.EmpresaID))

	if filtros.UsuarioID > 0 {
		query = query.Where(permiso.UsuarioIDEQ(filtros.UsuarioID))
	}
	if filtros.Estado != "" {
		query = query.Where(permiso.EstadoEQ(permiso.Estado(filtros.Estado)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if filtros.Pagina > 0 && filtros.Limite > 0 {
		offset := (filtros.Pagina - 1) * filtros.Limite
		query = query.Offset(offset).Limit(filtros.Limite)
	}

	permisosEnt, err := query.Order(ent.Desc(permiso.FieldFecha)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var resultados []*domain.Permiso
	for _, p := range permisosEnt {
		resultados = append(resultados, mapPermisoEntToDomain(p))
	}

	return resultados, total, nil
}

func (r *PermisoRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Permiso, error) {
	p, err := r.client.Permiso.Query().
		Where(
			permiso.IDEQ(id),
			permiso.EmpresaIDEQ(empresaID),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("permiso no encontrado")
		}
		return nil, fmt.Errorf("error al buscar permiso: %v", err)
	}

	return mapPermisoEntToDomain(p), nil
}

func (r *PermisoRepoEnt) Crear(ctx context.Context, p *domain.Permiso) (*domain.Permiso, error) {
	nuevo, err := r.client.Permiso.Create().
		SetEmpresaID(p.EmpresaID).
		SetUsuarioID(p.UsuarioID).
		SetFecha(p.Fecha).
		SetMotivo(p.Motivo).
		SetEstado(permiso.Estado(p.Estado)).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("error al crear permiso: %v", err)
	}

	return mapPermisoEntToDomain(nuevo), nil
}

func (r *PermisoRepoEnt) Actualizar(ctx context.Context, p *domain.Permiso) (*domain.Permiso, error) {
	actualizador := r.client.Permiso.UpdateOneID(p.ID).
		SetEstado(permiso.Estado(p.Estado))

	if p.Respuesta != nil {
		actualizador.SetRespuesta(*p.Respuesta)
	}

	actualizado, err := actualizador.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar permiso: %v", err)
	}

	return mapPermisoEntToDomain(actualizado), nil
}

func mapPermisoEntToDomain(p *ent.Permiso) *domain.Permiso {
	if p == nil {
		return nil
	}
	return &domain.Permiso{
		ID:        p.ID,
		EmpresaID: p.EmpresaID,
		UsuarioID: p.UsuarioID,
		Fecha:     p.Fecha,
		Motivo:    p.Motivo,
		Estado:    string(p.Estado),
		Respuesta: p.Respuesta,
	}
}
