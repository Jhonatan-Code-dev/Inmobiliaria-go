package repository

import (
	"context"
	"rentals-go/ent"
	"rentals-go/ent/contrato"
	"rentals-go/ent/serviciomedicion"
	"rentals-go/internal/domain"
)

type ServicioMedicionRepoEnt struct {
	client *ent.Client
}

func NewServicioMedicionRepo(client *ent.Client) *ServicioMedicionRepoEnt {
	return &ServicioMedicionRepoEnt{client: client}
}

func (r *ServicioMedicionRepoEnt) Listar(ctx context.Context, filtros domain.ServicioMedicionFiltros) ([]*domain.ServicioMedicion, int, error) {
	query := r.client.ServicioMedicion.Query()

	if filtros.ContratoID > 0 {
		query = query.Where(serviciomedicion.ContratoID(filtros.ContratoID))
	} else {
		// Si no hay contrato_id, filtramos por empresa vía Unidad -> Propiedad -> Empresa
		// Pero es más fácil si tenemos empresa_id en ServicioMedicion o si filtramos por el contrato activo de la empresa.
		// Por ahora filtramos por el contratoID si viene, si no, mostramos todos los de la empresa vinculada.
		query = query.Where(serviciomedicion.HasContratoWith(contrato.EmpresaIDEQ(filtros.EmpresaID)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	entities, err := query.Limit(filtros.PorPagina).Offset(offset).Order(ent.Desc(serviciomedicion.FieldFechaLectura)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.ServicioMedicion
	for _, e := range entities {
		results = append(results, mapServicioToDomain(e))
	}

	return results, total, nil
}

func (r *ServicioMedicionRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.ServicioMedicion, error) {
	e, err := r.client.ServicioMedicion.Query().
		Where(serviciomedicion.IDEQ(id), serviciomedicion.HasContratoWith(contrato.EmpresaIDEQ(empresaID))).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapServicioToDomain(e), nil
}

func (r *ServicioMedicionRepoEnt) Crear(ctx context.Context, s *domain.ServicioMedicion) (*domain.ServicioMedicion, error) {
	builder := r.client.ServicioMedicion.Create().
		SetUnidadID(s.UnidadID).
		SetTipoServicio(serviciomedicion.TipoServicio(s.TipoServicio)).
		SetLecturaAnterior(s.LecturaAnterior).
		SetLecturaActual(s.LecturaActual).
		SetConsumo(s.Consumo).
		SetMonto(int64(s.Monto * 100)).
		SetFechaLectura(s.FechaLectura).
		SetProcesado(s.Procesado)

	if s.ContratoID > 0 {
		builder.SetContratoID(s.ContratoID)
	}

	entServicio, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapServicioToDomain(entServicio), nil
}

func (r *ServicioMedicionRepoEnt) Actualizar(ctx context.Context, s *domain.ServicioMedicion) (*domain.ServicioMedicion, error) {
	builder := r.client.ServicioMedicion.UpdateOneID(s.ID).
		SetLecturaActual(s.LecturaActual).
		SetConsumo(s.Consumo).
		SetMonto(int64(s.Monto * 100)).
		SetProcesado(s.Procesado)

	entServicio, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapServicioToDomain(entServicio), nil
}

func (r *ServicioMedicionRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := r.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return r.client.ServicioMedicion.DeleteOneID(id).Exec(ctx)
}

func (r *ServicioMedicionRepoEnt) ObtenerUltimaLectura(ctx context.Context, contratoID int, tipo string) (*domain.ServicioMedicion, error) {
	e, err := r.client.ServicioMedicion.Query().
		Where(serviciomedicion.ContratoID(contratoID), serviciomedicion.TipoServicioEQ(serviciomedicion.TipoServicio(tipo))).
		Order(ent.Desc(serviciomedicion.FieldFechaLectura)).
		First(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return mapServicioToDomain(e), nil
}

func mapServicioToDomain(e *ent.ServicioMedicion) *domain.ServicioMedicion {
	cid := 0
	if e.ContratoID != nil {
		cid = *e.ContratoID
	}
	return &domain.ServicioMedicion{
		ID:              e.ID,
		UnidadID:        e.UnidadID,
		ContratoID:      cid,
		TipoServicio:    string(e.TipoServicio),
		LecturaAnterior: e.LecturaAnterior,
		LecturaActual:   e.LecturaActual,
		Consumo:         e.Consumo,
		Monto:           float64(e.Monto) / 100,
		FechaLectura:    e.FechaLectura,
		Procesado:       e.Procesado,
	}
}
