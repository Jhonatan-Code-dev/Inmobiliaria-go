package repository

import (
	"context"
	"rentals-go/ent"
	"rentals-go/ent/cita"
	"rentals-go/internal/domain"
)

type CitaRepoEnt struct {
	client *ent.Client
}

func NewCitaRepo(client *ent.Client) *CitaRepoEnt {
	return &CitaRepoEnt{client: client}
}

func (r *CitaRepoEnt) Listar(ctx context.Context, filtros domain.CitaFiltros) ([]*domain.Cita, int, error) {
	query := r.client.Cita.Query().Where(cita.EmpresaIDEQ(filtros.EmpresaID))

	if filtros.PropiedadID > 0 {
		query = query.Where(cita.PropiedadIDEQ(filtros.PropiedadID))
	}
	if filtros.UnidadID > 0 {
		query = query.Where(cita.UnidadIDEQ(filtros.UnidadID))
	}
	if filtros.Estado != "" {
		query = query.Where(cita.EstadoEQ(cita.Estado(filtros.Estado)))
	}
	if filtros.Busqueda != "" {
		query = query.Where(
			cita.Or(
				cita.NombreProspectoContainsFold(filtros.Busqueda),
				cita.TelefonoProspectoContainsFold(filtros.Busqueda),
				cita.ComentariosContainsFold(filtros.Busqueda),
			),
		)
	}
	if filtros.Desde != nil {
		query = query.Where(cita.FechaVisitaGTE(*filtros.Desde))
	}
	if filtros.Hasta != nil {
		query = query.Where(cita.FechaVisitaLTE(*filtros.Hasta))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Si no hay límites o paginación explícita (por ejemplo para el calendario mensual), no paginamos
	if filtros.Pagina > 0 && filtros.PorPagina > 0 {
		offset := (filtros.Pagina - 1) * filtros.PorPagina
		query = query.Limit(filtros.PorPagina).Offset(offset)
	}

	// Ordenamos las visitas de la más cercana a la más lejana
	entities, err := query.
		WithPropiedad().
		WithUnidad().
		WithCliente().
		Order(ent.Asc(cita.FieldFechaVisita)).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.Cita
	for _, e := range entities {
		results = append(results, mapCitaToDomain(e))
	}

	return results, total, nil
}

func (r *CitaRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Cita, error) {
	e, err := r.client.Cita.Query().
		Where(cita.IDEQ(id), cita.EmpresaIDEQ(empresaID)).
		WithPropiedad().
		WithUnidad().
		WithCliente().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapCitaToDomain(e), nil
}

func (r *CitaRepoEnt) Crear(ctx context.Context, c *domain.Cita) (*domain.Cita, error) {
	creator := r.client.Cita.Create().
		SetEmpresaID(c.EmpresaID).
		SetNombreProspecto(c.NombreProspecto).
		SetTelefonoProspecto(c.TelefonoProspecto).
		SetNillableCorreoProspecto(c.CorreoProspecto).
		SetFechaVisita(c.FechaVisita).
		SetEstado(cita.Estado(c.Estado)).
		SetNillableComentarios(c.Comentarios)

	if c.PropiedadID != nil {
		creator.SetPropiedadID(*c.PropiedadID)
	}
	if c.UnidadID != nil {
		creator.SetUnidadID(*c.UnidadID)
	}
	if c.ClienteID != nil {
		creator.SetClienteID(*c.ClienteID)
	}

	entCita, err := creator.Save(ctx)
	if err != nil {
		return nil, err
	}

	// Volvemos a consultar para tener cargadas las relaciones
	return r.BuscarPorID(ctx, entCita.ID, c.EmpresaID)
}

func (r *CitaRepoEnt) Actualizar(ctx context.Context, c *domain.Cita) (*domain.Cita, error) {
	updater := r.client.Cita.UpdateOneID(c.ID).
		SetNombreProspecto(c.NombreProspecto).
		SetTelefonoProspecto(c.TelefonoProspecto).
		SetNillableCorreoProspecto(c.CorreoProspecto).
		SetFechaVisita(c.FechaVisita).
		SetEstado(cita.Estado(c.Estado)).
		SetNillableComentarios(c.Comentarios)

	if c.PropiedadID != nil {
		updater.SetPropiedadID(*c.PropiedadID)
	} else {
		updater.ClearPropiedad()
	}

	if c.UnidadID != nil {
		updater.SetUnidadID(*c.UnidadID)
	} else {
		updater.ClearUnidad()
	}

	if c.ClienteID != nil {
		updater.SetClienteID(*c.ClienteID)
	} else {
		updater.ClearCliente()
	}

	entCita, err := updater.Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.BuscarPorID(ctx, entCita.ID, c.EmpresaID)
}

func (r *CitaRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := r.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return r.client.Cita.DeleteOneID(id).Exec(ctx)
}

func mapCitaToDomain(e *ent.Cita) *domain.Cita {
	c := &domain.Cita{
		ID:                e.ID,
		EmpresaID:         e.EmpresaID,
		PropiedadID:       e.PropiedadID,
		UnidadID:          e.UnidadID,
		ClienteID:         e.ClienteID,
		NombreProspecto:   e.NombreProspecto,
		TelefonoProspecto: e.TelefonoProspecto,
		CorreoProspecto:   e.CorreoProspecto,
		FechaVisita:       e.FechaVisita,
		Estado:            string(e.Estado),
		Comentarios:       e.Comentarios,
		CreadoEn:          e.CreadoEn,
	}

	if e.Edges.Propiedad != nil {
		c.PropiedadNombre = e.Edges.Propiedad.Nombre
	}
	if e.Edges.Unidad != nil {
		if e.Edges.Unidad.Nombre != nil {
			c.UnidadNombre = *e.Edges.Unidad.Nombre
		} else {
			c.UnidadNombre = e.Edges.Unidad.Codigo
		}
	}
	if e.Edges.Cliente != nil {
		nombreCompleto := e.Edges.Cliente.Nombres
		if e.Edges.Cliente.Apellidos != nil {
			nombreCompleto += " " + *e.Edges.Cliente.Apellidos
		}
		c.ClienteNombre = nombreCompleto
	}

	return c
}
