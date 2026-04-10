package repository

import (
	"context"
	"rentals-go/ent"
	"rentals-go/ent/ticket"
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

	if filtros.UnidadID > 0 {
		query = query.Where(ticket.UnidadIDEQ(filtros.UnidadID))
	}
	if filtros.Estado != "" {
		query = query.Where(ticket.EstadoEQ(ticket.Estado(filtros.Estado)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (filtros.Pagina - 1) * filtros.PorPagina
	entities, err := query.Limit(filtros.PorPagina).Offset(offset).Order(ent.Desc(ticket.FieldID)).All(ctx)
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

func mapTicketToDomain(e *ent.Ticket) *domain.Ticket {
	return &domain.Ticket{
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
}
