package repository

import (
	"context"
	"time"

	"rentals-go/ent"
	entGasto "rentals-go/ent/gasto"
	"rentals-go/internal/domain"
)

type GastoRepoEnt struct {
	client *ent.Client
}

func NewGastoRepo(client *ent.Client) *GastoRepoEnt {
	return &GastoRepoEnt{client: client}
}

func (r *GastoRepoEnt) ListarPaginado(ctx context.Context, filtros domain.GastoFiltros) ([]*domain.Gasto, int, error) {
	query := r.client.Gasto.Query().Where(entGasto.EmpresaID(filtros.EmpresaID))

	// Aplicar filtros de fecha
	if filtros.Fecha != nil {
		inicio := time.Date(filtros.Fecha.Year(), filtros.Fecha.Month(), filtros.Fecha.Day(), 0, 0, 0, 0, filtros.Fecha.Location())
		fin := inicio.AddDate(0, 0, 1)
		query = query.Where(entGasto.FechaGTE(inicio), entGasto.FechaLT(fin))
	} else if filtros.Desde != nil && filtros.Hasta != nil {
		query = query.Where(entGasto.FechaGTE(*filtros.Desde), entGasto.FechaLTE(*filtros.Hasta))
	} else if filtros.Anio > 0 {
		inicio := time.Date(filtros.Anio, 1, 1, 0, 0, 0, 0, time.UTC)
		fin := inicio.AddDate(1, 0, 0)
		if filtros.Mes > 0 {
			inicio = time.Date(filtros.Anio, time.Month(filtros.Mes), 1, 0, 0, 0, 0, time.UTC)
			fin = inicio.AddDate(0, 1, 0)
		}
		query = query.Where(entGasto.FechaGTE(inicio), entGasto.FechaLT(fin))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.Limite
	if offset < 0 {
		offset = 0
	}

	list, err := query.
		Limit(filtros.Limite).
		Offset(offset).
		Order(ent.Desc(entGasto.FieldFecha), ent.Desc(entGasto.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Gasto, 0, len(list))
	for _, e := range list {
		out = append(out, mapGastoEntity(e))
	}
	return out, total, nil
}

func (r *GastoRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Gasto, error) {
	e, err := r.client.Gasto.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapGastoEntity(e), nil
}

func (r *GastoRepoEnt) Crear(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	builder := r.client.Gasto.Create().
		SetEmpresaID(g.EmpresaID).
		SetMonto(g.Monto).
		SetFecha(g.Fecha).
		SetTipoPagoID(g.TipoPagoID).
		SetDescripcion(g.Descripcion)

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapGastoEntity(e), nil
}

func (r *GastoRepoEnt) Actualizar(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	builder := r.client.Gasto.UpdateOneID(g.ID).
		SetMonto(g.Monto).
		SetFecha(g.Fecha).
		SetTipoPagoID(g.TipoPagoID).
		SetDescripcion(g.Descripcion)

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapGastoEntity(e), nil
}

func (r *GastoRepoEnt) Eliminar(ctx context.Context, id int) error {
	return r.client.Gasto.DeleteOneID(id).Exec(ctx)
}

func mapGastoEntity(e *ent.Gasto) *domain.Gasto {
	return &domain.Gasto{
		ID:          e.ID,
		EmpresaID:   e.EmpresaID,
		Monto:       e.Monto,
		Fecha:       e.Fecha,
		TipoPagoID:  e.TipoPagoID,
		Descripcion: e.Descripcion,
	}
}
