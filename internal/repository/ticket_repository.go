package repository

import (
	"context"
	"rentals-go/ent"
	"rentals-go/ent/ticket"
	"rentals-go/ent/unidad"
	"rentals-go/internal/domain"
)

type TicketRepoEnt struct {
	client *ent.Client
}

func NewTicketRepo(client *ent.Client) *TicketRepoEnt {
	return &TicketRepoEnt{client: client}
}

func (r *TicketRepoEnt) Listar(ctx context.Context, filtros domain.TicketFiltros) ([]*domain.Ticket, int, error) {
	query := r.client.Ticket.Query().Where(ticket.EmpresaIDEQ(filtros.EmpresaID))

	if filtros.PropiedadID > 0 {
		query = query.Where(ticket.HasUnidadWith(unidad.PropiedadIDEQ(filtros.PropiedadID)))
	}
	if filtros.UnidadID > 0 {
		query = query.Where(ticket.UnidadIDEQ(filtros.UnidadID))
	}
	if filtros.Estado != "" {
		query = query.Where(ticket.EstadoEQ(ticket.Estado(filtros.Estado)))
	}
	if filtros.Busqueda != "" {
		query = query.Where(
			ticket.Or(
				ticket.AsuntoContainsFold(filtros.Busqueda),
				ticket.DescripcionContainsFold(filtros.Busqueda),
			),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	entities, err := query.Limit(filtros.PorPagina).Offset(offset).
		WithUnidad().
		WithCliente().
		Order(ent.Desc(ticket.FieldID)).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []*domain.Ticket
	for _, e := range entities {
		results = append(results, mapTicketToDomain(e))
	}

	return results, total, nil
}

func (r *TicketRepoEnt) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Ticket, error) {
	e, err := r.client.Ticket.Query().
		Where(ticket.IDEQ(id), ticket.EmpresaIDEQ(empresaID)).
		WithUnidad().
		WithCliente().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return mapTicketToDomain(e), nil
}

func (r *TicketRepoEnt) Crear(ctx context.Context, t *domain.Ticket) (*domain.Ticket, error) {
	entTicket, err := r.client.Ticket.Create().
		SetEmpresaID(t.EmpresaID).
		SetUnidadID(t.UnidadID).
		SetNillableClienteID(t.ClienteID).
		SetAsunto(t.Asunto).
		SetDescripcion(t.Descripcion).
		SetPrioridad(ticket.Prioridad(t.Prioridad)).
		SetEstado(ticket.Estado(t.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapTicketToDomain(entTicket), nil
}

func (r *TicketRepoEnt) Actualizar(ctx context.Context, t *domain.Ticket) (*domain.Ticket, error) {
	entTicket, err := r.client.Ticket.UpdateOneID(t.ID).
		SetAsunto(t.Asunto).
		SetDescripcion(t.Descripcion).
		SetPrioridad(ticket.Prioridad(t.Prioridad)).
		SetEstado(ticket.Estado(t.Estado)).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return mapTicketToDomain(entTicket), nil
}

func (r *TicketRepoEnt) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := r.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return r.client.Ticket.DeleteOneID(id).Exec(ctx)
}

func (r *TicketRepoEnt) ObtenerResumen(ctx context.Context, empresaID int, propiedadID int) (*domain.TicketResumen, error) {
	query := r.client.Ticket.Query().Where(ticket.EmpresaIDEQ(empresaID))
	
	if propiedadID > 0 {
		query = query.Where(ticket.HasUnidadWith(unidad.PropiedadIDEQ(propiedadID)))
	}

	resumen := &domain.TicketResumen{}

	// Se podría hacer con GroupBy, pero para pocos estados múltiples count es aceptable y claro
	var err error
	resumen.Total, err = query.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	
	resumen.Abiertos, _ = query.Clone().Where(ticket.EstadoEQ(ticket.EstadoAbierto)).Count(ctx)
	resumen.EnProgreso, _ = query.Clone().Where(ticket.EstadoEQ(ticket.EstadoEnProgreso)).Count(ctx)
	resumen.Resueltos, _ = query.Clone().Where(ticket.EstadoEQ(ticket.EstadoResuelto)).Count(ctx)
	resumen.Cerrados, _ = query.Clone().Where(ticket.EstadoEQ(ticket.EstadoCerrado)).Count(ctx)

	return resumen, nil
}

func (r *TicketRepoEnt) ListarColaTrabajo(ctx context.Context, filtros domain.TicketFiltros) ([]*domain.Ticket, int, error) {
	// Solo tickets que requieren atención (abierto o en progreso)
	query := r.client.Ticket.Query().
		Where(
			ticket.EmpresaID(filtros.EmpresaID),
			ticket.EstadoIn(ticket.EstadoAbierto, ticket.EstadoEnProgreso),
		)

	if filtros.Busqueda != "" {
		query = query.Where(
			ticket.Or(
				ticket.AsuntoContainsFold(filtros.Busqueda),
				ticket.DescripcionContainsFold(filtros.Busqueda),
			),
		)
	}

	// Filtros opcionales de inmueble/unidad
	if filtros.PropiedadID > 0 {
		query = query.Where(ticket.HasUnidadWith(unidad.PropiedadIDEQ(filtros.PropiedadID)))
	}
	if filtros.UnidadID > 0 {
		query = query.Where(ticket.UnidadID(filtros.UnidadID))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	if offset < 0 {
		offset = 0
	}

	// Orden: Prioridad (Alta primero) y luego por antigüedad (el más viejo primero)
	list, err := query.
		WithUnidad().
		WithCliente().
		Limit(filtros.PorPagina).
		Offset(offset).
		Order(
			ent.Desc(ticket.FieldPrioridad), // alta > media > baja
			ent.Asc(ticket.FieldCreadoEn),    // Antiguo primero
		).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	out := make([]*domain.Ticket, 0, len(list))
	for _, item := range list {
		out = append(out, mapTicketToDomain(item))
	}
	return out, total, nil
}

func mapTicketToDomain(e *ent.Ticket) *domain.Ticket {
	t := &domain.Ticket{
		ID:            e.ID,
		EmpresaID:     e.EmpresaID,
		UnidadID:      e.UnidadID,
		ClienteID:     e.ClienteID,
		Asunto:        e.Asunto,
		Descripcion:   e.Descripcion,
		Prioridad:     string(e.Prioridad),
		Estado:        string(e.Estado),
		FechaApertura: e.CreadoEn,
	}
	if e.Edges.Unidad != nil {
		if e.Edges.Unidad.Nombre != nil {
			t.UnidadNombre = *e.Edges.Unidad.Nombre
		} else {
			t.UnidadNombre = e.Edges.Unidad.Codigo
		}
	}
	if e.Edges.Cliente != nil {
		nombreCompleto := e.Edges.Cliente.Nombres
		if e.Edges.Cliente.Apellidos != nil {
			nombreCompleto += " " + *e.Edges.Cliente.Apellidos
		}
		t.ClienteNombre = nombreCompleto
	}
	return t
}
