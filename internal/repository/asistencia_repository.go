package repository

import (
	"context"
	"fmt"
	"time"

	"rentals-go/ent"
	"rentals-go/ent/asistencia"
	"rentals-go/internal/domain"
)

type AsistenciaRepoEnt struct {
	client *ent.Client
}

func NewAsistenciaRepo(client *ent.Client) *AsistenciaRepoEnt {
	return &AsistenciaRepoEnt{client: client}
}

func (r *AsistenciaRepoEnt) ListarPaginado(ctx context.Context, filtros domain.AsistenciaFiltros) ([]*domain.Asistencia, int, error) {
	query := r.client.Asistencia.Query().
		Where(asistencia.EmpresaIDEQ(filtros.EmpresaID))

	if filtros.UsuarioID > 0 {
		query = query.Where(asistencia.UsuarioIDEQ(filtros.UsuarioID))
	}
	if filtros.Estado != "" {
		query = query.Where(asistencia.EstadoEQ(asistencia.Estado(filtros.Estado)))
	}
	if filtros.Desde != nil {
		query = query.Where(asistencia.FechaGTE(*filtros.Desde))
	}
	if filtros.Hasta != nil {
		query = query.Where(asistencia.FechaLTE(*filtros.Hasta))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if filtros.Pagina > 0 && filtros.Limite > 0 {
		offset := (filtros.Pagina - 1) * filtros.Limite
		query = query.Offset(offset).Limit(filtros.Limite)
	}

	// Cargar información del usuario y ordenar por fecha descendente
	asistenciasEnt, err := query.
		WithUsuario().
		Order(ent.Desc(asistencia.FieldFecha)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var resultados []*domain.Asistencia
	for _, a := range asistenciasEnt {
		resultados = append(resultados, mapAsistenciaEntToDomain(a))
	}

	return resultados, total, nil
}

func (r *AsistenciaRepoEnt) BuscarPorFechaUsuario(ctx context.Context, usuarioID int, empresaID int, fecha time.Time) (*domain.Asistencia, error) {
	// Normalizar fecha para buscar solo por día
	inicioDia := time.Date(fecha.Year(), fecha.Month(), fecha.Day(), 0, 0, 0, 0, fecha.Location())
	finDia := inicioDia.Add(24 * time.Hour).Add(-time.Nanosecond)

	a, err := r.client.Asistencia.Query().
		Where(
			asistencia.UsuarioIDEQ(usuarioID),
			asistencia.EmpresaIDEQ(empresaID),
			asistencia.FechaGTE(inicioDia),
			asistencia.FechaLTE(finDia),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil // No se considera un error crítico si no ha marcado aún
		}
		return nil, fmt.Errorf("error al buscar asistencia: %v", err)
	}

	return mapAsistenciaEntToDomain(a), nil
}

func (r *AsistenciaRepoEnt) Crear(ctx context.Context, a *domain.Asistencia) (*domain.Asistencia, error) {
	creador := r.client.Asistencia.Create().
		SetEmpresaID(a.EmpresaID).
		SetUsuarioID(a.UsuarioID).
		SetFecha(a.Fecha).
		SetEstado(asistencia.Estado(a.Estado))

	if a.HoraEntrada != nil {
		creador.SetHoraEntrada(*a.HoraEntrada)
	}
	if a.HoraSalida != nil {
		creador.SetHoraSalida(*a.HoraSalida)
	}
	if a.Notas != nil {
		creador.SetNotas(*a.Notas)
	}
	if a.HorasTrabajadas != nil {
		creador.SetHorasTrabajadas(*a.HorasTrabajadas)
	}

	nuevo, err := creador.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al crear asistencia: %v", err)
	}

	return mapAsistenciaEntToDomain(nuevo), nil
}

func (r *AsistenciaRepoEnt) Actualizar(ctx context.Context, a *domain.Asistencia) (*domain.Asistencia, error) {
	actualizador := r.client.Asistencia.UpdateOneID(a.ID).
		SetEstado(asistencia.Estado(a.Estado))

	if a.HoraEntrada != nil {
		actualizador.SetHoraEntrada(*a.HoraEntrada)
	}
	if a.HoraSalida != nil {
		actualizador.SetHoraSalida(*a.HoraSalida)
	}
	if a.Notas != nil {
		actualizador.SetNotas(*a.Notas)
	}
	if a.HorasTrabajadas != nil {
		actualizador.SetHorasTrabajadas(*a.HorasTrabajadas)
	}

	actualizado, err := actualizador.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar asistencia: %v", err)
	}

	return mapAsistenciaEntToDomain(actualizado), nil
}

func (r *AsistenciaRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	return r.client.Asistencia.DeleteOneID(id).
		Where(asistencia.EmpresaIDEQ(empresaID)).
		Exec(ctx)
}

func mapAsistenciaEntToDomain(a *ent.Asistencia) *domain.Asistencia {
	if a == nil {
		return nil
	}
	d := &domain.Asistencia{
		ID:              a.ID,
		EmpresaID:       a.EmpresaID,
		UsuarioID:       a.UsuarioID,
		Fecha:           a.Fecha,
		HoraEntrada:     a.HoraEntrada,
		HoraSalida:      a.HoraSalida,
		Estado:          string(a.Estado),
		Notas:           a.Notas,
		HorasTrabajadas: a.HorasTrabajadas,
	}

	if a.Edges.Usuario != nil {
		d.UsuarioNombre = a.Edges.Usuario.Usuario
	}

	return d
}
