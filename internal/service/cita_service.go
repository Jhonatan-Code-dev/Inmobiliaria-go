package service

import (
	"context"
	"errors"
	"rentals-go/internal/domain"
)

type CitaService struct {
	repo domain.CitaRepository
}

func NewCitaService(repo domain.CitaRepository) *CitaService {
	return &CitaService{repo: repo}
}

func (s *CitaService) Listar(ctx context.Context, filtros domain.CitaFiltros) ([]*domain.Cita, int, error) {
	// Permitimos consulta de todas las citas del mes sin paginación si pagina o porPagina son 0
	if filtros.Pagina > 0 && filtros.PorPagina <= 0 {
		filtros.PorPagina = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *CitaService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Cita, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *CitaService) Crear(ctx context.Context, req *domain.RegistroCita, empresaID int) (*domain.Cita, error) {
	if req.NombreProspecto == "" {
		return nil, errors.New("el nombre del prospecto es requerido")
	}
	if req.TelefonoProspecto == "" {
		return nil, errors.New("el teléfono del prospecto es requerido")
	}
	if req.FechaVisita.IsZero() {
		return nil, errors.New("la fecha y hora de la visita es requerida")
	}
	if req.Estado == "" {
		req.Estado = "programada"
	}

	c := &domain.Cita{
		EmpresaID:         empresaID,
		PropiedadID:       req.PropiedadID,
		UnidadID:          req.UnidadID,
		ClienteID:         req.ClienteID,
		NombreProspecto:   req.NombreProspecto,
		TelefonoProspecto: req.TelefonoProspecto,
		CorreoProspecto:   req.CorreoProspecto,
		FechaVisita:       req.FechaVisita,
		Estado:            req.Estado,
		Comentarios:       req.Comentarios,
	}

	return s.repo.Crear(ctx, c)
}

func (s *CitaService) Actualizar(ctx context.Context, id int, empresaID int, req *domain.RegistroCita) (*domain.Cita, error) {
	c, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}

	if req.NombreProspecto != "" {
		c.NombreProspecto = req.NombreProspecto
	}
	if req.TelefonoProspecto != "" {
		c.TelefonoProspecto = req.TelefonoProspecto
	}
	if req.CorreoProspecto != nil {
		c.CorreoProspecto = req.CorreoProspecto
	}
	if !req.FechaVisita.IsZero() {
		c.FechaVisita = req.FechaVisita
	}
	if req.Estado != "" {
		c.Estado = req.Estado
	}
	if req.Comentarios != nil {
		c.Comentarios = req.Comentarios
	}

	// Actualizamos relaciones opcionales
	c.PropiedadID = req.PropiedadID
	c.UnidadID = req.UnidadID
	c.ClienteID = req.ClienteID

	return s.repo.Actualizar(ctx, c)
}

func (s *CitaService) CambiarEstado(ctx context.Context, id int, empresaID int, estado string) (*domain.Cita, error) {
	c, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}

	c.Estado = estado
	return s.repo.Actualizar(ctx, c)
}

func (s *CitaService) Eliminar(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}
