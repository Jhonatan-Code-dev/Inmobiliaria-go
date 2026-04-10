package repository

import (
	"context"
	"rentals-go/ent"
	"rentals-go/ent/cargo"
	"rentals-go/ent/contrato"
	"rentals-go/internal/domain"
)

type CargoRepoEnt struct {
	client *ent.Client
}

func NewCargoRepo(client *ent.Client) *CargoRepoEnt {
	return &CargoRepoEnt{client: client}
}

func (r *CargoRepoEnt) Listar(ctx context.Context, filtros domain.CargoFiltros) ([]*domain.Cargo, int, error) {
	query := r.client.Cargo.Query().
		Where(cargo.HasContratoWith(contrato.EmpresaID(filtros.EmpresaID)))

	if filtros.ContratoID > 0 {
		query = query.Where(cargo.ContratoID(filtros.ContratoID))
	}
	if filtros.Estado != "" {
		query = query.Where(cargo.EstadoEQ(cargo.Estado(filtros.Estado)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	entities, err := query.Limit(filtros.PorPagina).Offset(offset).Order(ent.Desc(cargo.FieldFechaVencimiento)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.Cargo
	for _, e := range entities {
		results = append(results, mapCargoToDomain(e))
	}

	return results, total, nil
}

func (r *CargoRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Cargo, error) {
	e, err := r.client.Cargo.Query().
		Where(cargo.IDEQ(id), cargo.HasContratoWith(contrato.EmpresaIDEQ(empresaID))).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapCargoToDomain(e), nil
}

func (r *CargoRepoEnt) Crear(ctx context.Context, c *domain.Cargo) (*domain.Cargo, error) {
	entCargo, err := r.client.Cargo.Create().
		SetContratoID(c.ContratoID).
		SetConcepto(cargo.Concepto(c.Concepto)).
		SetDescripcion(c.Descripcion).
		SetMoneda(c.Moneda).
		SetFechaVencimiento(c.FechaVencimiento).
		SetMonto(int64(c.Monto * 100)).
		SetSaldo(int64(c.Monto * 100)).
		SetEstado(cargo.EstadoPendiente).
		SetGeneradoAutomaticamente(c.GeneradoAutomaticamente).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapCargoToDomain(entCargo), nil
}

func (r *CargoRepoEnt) Actualizar(ctx context.Context, c *domain.Cargo) (*domain.Cargo, error) {
	entCargo, err := r.client.Cargo.UpdateOneID(c.ID).
		SetConcepto(cargo.Concepto(c.Concepto)).
		SetDescripcion(c.Descripcion).
		SetFechaVencimiento(c.FechaVencimiento).
		SetMonto(int64(c.Monto * 100)).
		SetSaldo(int64(c.Saldo * 100)).
		SetEstado(cargo.Estado(c.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapCargoToDomain(entCargo), nil
}

func (r *CargoRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := r.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return r.client.Cargo.DeleteOneID(id).Exec(ctx)
}

func mapCargoToDomain(e *ent.Cargo) *domain.Cargo {
	return &domain.Cargo{
		ID:                      e.ID,
		ContratoID:              e.ContratoID,
		Concepto:                string(e.Concepto),
		Descripcion:             ptrToString(e.Descripcion),
		Moneda:                  e.Moneda,
		PeriodoInicio:           e.PeriodoInicio,
		PeriodoFin:              e.PeriodoFin,
		FechaEmision:            e.FechaEmision,
		FechaVencimiento:        e.FechaVencimiento,
		Monto:                   float64(e.Monto) / 100,
		Saldo:                   float64(e.Saldo) / 100,
		Estado:                  string(e.Estado),
		GeneradoAutomaticamente: e.GeneradoAutomaticamente,
	}
}
