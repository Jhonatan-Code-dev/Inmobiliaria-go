package repository

import (
	"context"

	"rentals-go/ent"
	entMov "rentals-go/ent/movimientocaja"
	"rentals-go/internal/domain"
)

type MovimientoCajaRepoEnt struct {
	client *ent.Client
}

func NewMovimientoCajaRepo(client *ent.Client) *MovimientoCajaRepoEnt {
	return &MovimientoCajaRepoEnt{client: client}
}

func (r *MovimientoCajaRepoEnt) Crear(ctx context.Context, m *domain.MovimientoCaja) (*domain.MovimientoCaja, error) {
	builder := r.client.MovimientoCaja.Create().
		SetEmpresaID(m.EmpresaID).
		SetNillablePagoID(m.PagoID).
		SetNillableGastoID(m.GastoID).
		SetTipo(entMov.Tipo(m.Tipo)).
		SetConcepto(m.Concepto).
		SetFechaMovimiento(m.FechaMovimiento).
		SetMoneda(m.Moneda).
		SetMonto(m.Monto).
		SetMetodo(entMov.Metodo(m.Metodo)).
		SetNillableReferencia(m.Referencia).
		SetNillableObservaciones(m.Observaciones)

	e, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapMovimientoEntity(e), nil
}

func mapMovimientoEntity(e *ent.MovimientoCaja) *domain.MovimientoCaja {
	return &domain.MovimientoCaja{
		ID:              e.ID,
		EmpresaID:       e.EmpresaID,
		PagoID:          e.PagoID,
		GastoID:         e.GastoID,
		Tipo:            string(e.Tipo),
		Concepto:        e.Concepto,
		FechaMovimiento: e.FechaMovimiento,
		Moneda:          e.Moneda,
		Monto:           e.Monto,
		Metodo:          string(e.Metodo),
		Referencia:      e.Referencia,
		Observaciones:   e.Observaciones,
	}
}
