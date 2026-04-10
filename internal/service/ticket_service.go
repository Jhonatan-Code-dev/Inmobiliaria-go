package service

import (
	"context"
	"rentals-go/internal/domain"
)

type TicketService struct {
	repo domain.TicketRepository
}

func NewTicketService(repo domain.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) Listar(ctx context.Context, filtros domain.TicketFiltros) ([]*domain.Ticket, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.PorPagina <= 0 {
		filtros.PorPagina = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *TicketService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Ticket, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *TicketService) Crear(ctx context.Context, r *domain.RegistroTicket, empresaID int) (*domain.Ticket, error) {
	t := &domain.Ticket{
		EmpresaID:   empresaID,
		UnidadID:    r.UnidadID,
		Asunto:      r.Asunto,
		Descripcion: r.Descripcion,
		Prioridad:   r.Prioridad,
		Estado:      "abierto",
	}
	return s.repo.Crear(ctx, t)
}

func (s *TicketService) Actualizar(ctx context.Context, id int, empresaID int, r *domain.RegistroTicket, estado string) (*domain.Ticket, error) {
	t, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}

	if r.Asunto != "" {
		t.Asunto = r.Asunto
	}
	if r.Descripcion != "" {
		t.Descripcion = r.Descripcion
	}
	if r.Prioridad != "" {
		t.Prioridad = r.Prioridad
	}
	if estado != "" {
		t.Estado = estado
	}

	return s.repo.Actualizar(ctx, t)
}

func (s *TicketService) Eliminar(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}
