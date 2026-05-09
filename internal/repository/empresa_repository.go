package repository

import (
	"context"

	"rentals-go/ent"
	entEmpresa "rentals-go/ent/empresa"
	"rentals-go/internal/domain"
)

type EmpresaRepoEnt struct {
	client *ent.Client
}

func NewEmpresaRepo(client *ent.Client) *EmpresaRepoEnt {
	return &EmpresaRepoEnt{client: client}
}

func (r *EmpresaRepoEnt) ListarPaginado(ctx context.Context, limite, offset int, busqueda string) ([]*domain.Empresa, int, error) {
	query := r.client.Empresa.Query()

	if busqueda != "" {
		query = query.Where(
			entEmpresa.NombreContains(busqueda),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	list, err := query.
		Limit(limite).
		Offset(offset).
		Order(ent.Desc(entEmpresa.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Empresa, 0, len(list))
	for _, e := range list {
		out = append(out, mapEmpresaEntity(e))
	}
	return out, total, nil
}

func (r *EmpresaRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Empresa, error) {
	e, err := r.client.Empresa.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Crear(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	e, err := r.client.Empresa.Create().
		SetNombre(emp.Nombre).
		SetNillablePais(nilIfEmpty(emp.Pais)).
		SetMoneda(emp.Moneda).
		SetMaximoUsuarios(defaultInt(emp.MaximoUsuarios, 1)).
		SetEstado(emp.Estado).
		SetNillableVencimiento(nilIfTimeZero(emp.Vencimiento)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Actualizar(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	e, err := r.client.Empresa.UpdateOneID(emp.ID).
		SetNombre(emp.Nombre).
		SetNillablePais(nilIfEmpty(emp.Pais)).
		SetMoneda(emp.Moneda).
		SetMaximoUsuarios(defaultInt(emp.MaximoUsuarios, 1)).
		SetEstado(emp.Estado).
		SetNillableVencimiento(nilIfTimeZero(emp.Vencimiento)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapEmpresaEntity(e), nil
}

func (r *EmpresaRepoEnt) Eliminar(ctx context.Context, id int) error {
	return r.client.Empresa.DeleteOneID(id).Exec(ctx)
}

func (r *EmpresaRepoEnt) ActualizarConfiguracionAsistencia(ctx context.Context, empresaID int, config *domain.ConfiguracionAsistencia) error {
	_, err := r.client.Empresa.UpdateOneID(empresaID).
		SetHorarioEntradaDefecto(config.HoraEntrada).
		SetHorarioSalidaDefecto(config.HoraSalida).
		SetToleranciaDefecto(config.ToleranciaMinutos).
		SetDiasLaborablesDefecto(config.DiasLaborables).
		Save(ctx)
	return err
}

func mapEmpresaEntity(e *ent.Empresa) *domain.Empresa {
	return &domain.Empresa{
		ID:                     e.ID,
		Nombre:                 e.Nombre,
		Pais:                   ptrToString(e.Pais),
		Moneda:                 e.Moneda,
		MaximoUsuarios:         e.MaximoUsuarios,
		Estado:                 e.Estado,
		Vencimiento:            ptrToTime(e.Vencimiento),
		CreadoEn:               e.CreadoEn,
		HorarioEntradaDefecto:  e.HorarioEntradaDefecto,
		HorarioSalidaDefecto:   e.HorarioSalidaDefecto,
		ToleranciaDefecto:      e.ToleranciaDefecto,
		DiasLaborablesDefecto: e.DiasLaborablesDefecto,
	}
}

func defaultInt(val, def int) int {
	if val <= 0 {
		return def
	}
	return val
}

func defaultStringValue(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
