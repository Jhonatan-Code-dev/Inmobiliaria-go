package repository

import (
	"context"

	"rentals-go/ent"
	entPropiedad "rentals-go/ent/propiedad"
	entUnidad "rentals-go/ent/unidad"
	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/money"
)

type InmuebleRepoEnt struct {
	client *ent.Client
}

func NewInmuebleRepo(client *ent.Client) *InmuebleRepoEnt {
	return &InmuebleRepoEnt{client: client}
}

func (r *InmuebleRepoEnt) ListarPaginado(ctx context.Context, filtros domain.InmuebleFiltros) ([]*domain.Inmueble, int, error) {
	query := r.client.Propiedad.Query().Where(entPropiedad.EmpresaID(filtros.EmpresaID))

	if filtros.Busqueda != "" {
		query = query.Where(
			entPropiedad.Or(
				entPropiedad.NombreContainsFold(filtros.Busqueda),
				entPropiedad.DireccionContainsFold(filtros.Busqueda),
				entPropiedad.CiudadContainsFold(filtros.Busqueda),
			),
		)
	}
	if filtros.Estado != "" {
		query = query.Where(entPropiedad.EstadoEQ(entPropiedad.Estado(filtros.Estado)))
	}
	if filtros.Tipo != "" {
		query = query.Where(entPropiedad.TipoEQ(entPropiedad.Tipo(filtros.Tipo)))
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
		Order(ent.Desc(entPropiedad.FieldCreadoEn), ent.Desc(entPropiedad.FieldID)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Inmueble, 0, len(list))
	for _, item := range list {
		out = append(out, mapPropiedadEntity(item, nil))
	}
	return out, total, nil
}

func (r *InmuebleRepoEnt) BuscarPorID(ctx context.Context, id int) (*domain.Inmueble, error) {
	item, err := r.client.Propiedad.Query().
		Where(entPropiedad.IDEQ(id)).
		WithUnidades().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	unidades := make([]*domain.Unidad, 0, len(item.Edges.Unidades))
	for _, unidad := range item.Edges.Unidades {
		unidades = append(unidades, mapUnidadEntity(unidad))
	}

	return mapPropiedadEntity(item, unidades), nil
}

func (r *InmuebleRepoEnt) Crear(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	item, err := r.client.Propiedad.Create().
		SetEmpresaID(inmueble.EmpresaID).
		SetNombre(inmueble.Nombre).
		SetTipo(entPropiedad.Tipo(inmueble.Tipo)).
		SetNillableDescripcion(inmueble.Descripcion).
		SetDireccion(inmueble.Direccion).
		SetNillableCiudad(inmueble.Ciudad).
		SetNillableRegion(inmueble.Region).
		SetNillablePais(inmueble.Pais).
		SetNillableCodigoPostal(inmueble.CodigoPostal).
		SetTotalPisos(inmueble.TotalPisos).
		SetTotalUnidades(inmueble.TotalUnidades).
		SetEstado(entPropiedad.Estado(inmueble.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapPropiedadEntity(item, nil), nil
}

func (r *InmuebleRepoEnt) Actualizar(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	item, err := r.client.Propiedad.UpdateOneID(inmueble.ID).
		SetNombre(inmueble.Nombre).
		SetTipo(entPropiedad.Tipo(inmueble.Tipo)).
		SetNillableDescripcion(inmueble.Descripcion).
		SetDireccion(inmueble.Direccion).
		SetNillableCiudad(inmueble.Ciudad).
		SetNillableRegion(inmueble.Region).
		SetNillablePais(inmueble.Pais).
		SetNillableCodigoPostal(inmueble.CodigoPostal).
		SetTotalPisos(inmueble.TotalPisos).
		SetTotalUnidades(inmueble.TotalUnidades).
		SetEstado(entPropiedad.Estado(inmueble.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return mapPropiedadEntity(item, nil), nil
}

func (r *InmuebleRepoEnt) Eliminar(ctx context.Context, id int) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer rollbackTx(tx)

	if _, err := tx.Unidad.Delete().Where(entUnidad.PropiedadIDEQ(id)).Exec(ctx); err != nil {
		return err
	}
	if err := tx.Propiedad.DeleteOneID(id).Exec(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *InmuebleRepoEnt) ListarUnidades(ctx context.Context, propiedadID int) ([]*domain.Unidad, error) {
	propiedad, err := r.client.Propiedad.Query().
		Where(entPropiedad.IDEQ(propiedadID)).
		WithUnidades(func(q *ent.UnidadQuery) {
			q.Order(ent.Asc(entUnidad.FieldCodigo))
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*domain.Unidad, 0, len(propiedad.Edges.Unidades))
	for _, item := range propiedad.Edges.Unidades {
		out = append(out, mapUnidadEntity(item))
	}
	return out, nil
}

func (r *InmuebleRepoEnt) BuscarUnidadPorID(ctx context.Context, id int) (*domain.Unidad, error) {
	item, err := r.client.Unidad.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapUnidadEntity(item), nil
}

func (r *InmuebleRepoEnt) CrearUnidad(ctx context.Context, unidad *domain.Unidad) (*domain.Unidad, error) {
	item, err := r.client.Unidad.Create().
		SetPropiedadID(unidad.PropiedadID).
		SetCodigo(unidad.Codigo).
		SetNillableNombre(unidad.Nombre).
		SetTipo(entUnidad.Tipo(unidad.Tipo)).
		SetNillableNumeroPiso(unidad.NumeroPiso).
		SetDormitorios(unidad.Dormitorios).
		SetBanos(unidad.Banos).
		SetNillableAreaM2(unidad.AreaM2).
		SetCapacidad(unidad.Capacidad).
		SetMoneda(unidad.Moneda).
		SetPrecioBase(unidad.PrecioBaseCents).
		SetDepositoRequerido(unidad.DepositoReqCents).
		SetIncluyeAgua(unidad.IncluyeAgua).
		SetIncluyeLuz(unidad.IncluyeLuz).
		SetIncluyeInternet(unidad.IncluyeInternet).
		SetNillableNotas(unidad.Notas).
		SetEstado(entUnidad.Estado(unidad.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.sincronizarTotalUnidades(ctx, unidad.PropiedadID); err != nil {
		return nil, err
	}
	return mapUnidadEntity(item), nil
}

func (r *InmuebleRepoEnt) ActualizarUnidad(ctx context.Context, unidad *domain.Unidad) (*domain.Unidad, error) {
	item, err := r.client.Unidad.UpdateOneID(unidad.ID).
		SetCodigo(unidad.Codigo).
		SetNillableNombre(unidad.Nombre).
		SetTipo(entUnidad.Tipo(unidad.Tipo)).
		SetNillableNumeroPiso(unidad.NumeroPiso).
		SetDormitorios(unidad.Dormitorios).
		SetBanos(unidad.Banos).
		SetNillableAreaM2(unidad.AreaM2).
		SetCapacidad(unidad.Capacidad).
		SetMoneda(unidad.Moneda).
		SetPrecioBase(unidad.PrecioBaseCents).
		SetDepositoRequerido(unidad.DepositoReqCents).
		SetIncluyeAgua(unidad.IncluyeAgua).
		SetIncluyeLuz(unidad.IncluyeLuz).
		SetIncluyeInternet(unidad.IncluyeInternet).
		SetNillableNotas(unidad.Notas).
		SetEstado(entUnidad.Estado(unidad.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapUnidadEntity(item), nil
}

func (r *InmuebleRepoEnt) EliminarUnidad(ctx context.Context, id int) error {
	item, err := r.client.Unidad.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := r.client.Unidad.DeleteOneID(id).Exec(ctx); err != nil {
		return err
	}
	return r.sincronizarTotalUnidades(ctx, item.PropiedadID)
}

func mapPropiedadEntity(item *ent.Propiedad, unidades []*domain.Unidad) *domain.Inmueble {
	return &domain.Inmueble{
		ID:            item.ID,
		EmpresaID:     item.EmpresaID,
		Nombre:        item.Nombre,
		Tipo:          string(item.Tipo),
		Descripcion:   item.Descripcion,
		Direccion:     item.Direccion,
		Ciudad:        item.Ciudad,
		Region:        item.Region,
		Pais:          item.Pais,
		CodigoPostal:  item.CodigoPostal,
		TotalPisos:    item.TotalPisos,
		TotalUnidades: item.TotalUnidades,
		Estado:        string(item.Estado),
		CreadoEn:      item.CreadoEn,
		Unidades:      unidades,
	}
}

func mapUnidadEntity(item *ent.Unidad) *domain.Unidad {
	return &domain.Unidad{
		ID:                item.ID,
		PropiedadID:       item.PropiedadID,
		Codigo:            item.Codigo,
		Nombre:            item.Nombre,
		Tipo:              string(item.Tipo),
		NumeroPiso:        item.NumeroPiso,
		Dormitorios:       item.Dormitorios,
		Banos:             item.Banos,
		AreaM2:            item.AreaM2,
		Capacidad:         item.Capacidad,
		Moneda:            item.Moneda,
		PrecioBase:        money.NewAmountFromCents(item.PrecioBase).Float64(),
		PrecioBaseCents:   item.PrecioBase,
		DepositoRequerido: money.NewAmountFromCents(item.DepositoRequerido).Float64(),
		DepositoReqCents:  item.DepositoRequerido,
		IncluyeAgua:       item.IncluyeAgua,
		IncluyeLuz:        item.IncluyeLuz,
		IncluyeInternet:   item.IncluyeInternet,
		Notas:             item.Notas,
		Estado:            string(item.Estado),
		CreadoEn:          item.CreadoEn,
	}
}

func (r *InmuebleRepoEnt) sincronizarTotalUnidades(ctx context.Context, propiedadID int) error {
	total, err := r.client.Unidad.Query().Where(entUnidad.PropiedadIDEQ(propiedadID)).Count(ctx)
	if err != nil {
		return err
	}
	_, err = r.client.Propiedad.UpdateOneID(propiedadID).SetTotalUnidades(total).Save(ctx)
	return err
}
