package service

import (
	"context"
	"fmt"

	"rentals-go/internal/domain"
)

type InmuebleService struct {
	repo domain.InmuebleRepository
}

func NewInmuebleService(repo domain.InmuebleRepository) *InmuebleService {
	return &InmuebleService{repo: repo}
}

func (s *InmuebleService) Listar(ctx context.Context, filtros domain.InmuebleFiltros) ([]*domain.Inmueble, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 || filtros.Limite > 10 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *InmuebleService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Inmueble, error) {
	inmueble, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if inmueble.EmpresaID != empresaID {
		return nil, fmt.Errorf("%w: inmueble no pertenece a la empresa", domain.ErrForbidden)
	}
	return inmueble, nil
}

func (s *InmuebleService) Crear(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	if inmueble.Tipo == "" {
		inmueble.Tipo = "casa"
	}
	if inmueble.Estado == "" {
		inmueble.Estado = "activa"
	}
	if inmueble.TotalPisos <= 0 {
		inmueble.TotalPisos = 1
	}
	if inmueble.TotalUnidades <= 0 {
		inmueble.TotalUnidades = 1
	}
	return s.repo.Crear(ctx, inmueble)
}

func (s *InmuebleService) Actualizar(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	actual, err := s.repo.BuscarPorID(ctx, inmueble.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if actual.EmpresaID != inmueble.EmpresaID {
		return nil, fmt.Errorf("%w: no autorizado para actualizar este inmueble", domain.ErrForbidden)
	}
	if inmueble.Tipo == "" {
		inmueble.Tipo = actual.Tipo
	}
	if inmueble.Estado == "" {
		inmueble.Estado = actual.Estado
	}
	if inmueble.TotalPisos <= 0 {
		inmueble.TotalPisos = actual.TotalPisos
	}
	if inmueble.TotalUnidades <= 0 {
		inmueble.TotalUnidades = actual.TotalUnidades
	}
	return s.repo.Actualizar(ctx, inmueble)
}

func (s *InmuebleService) Eliminar(ctx context.Context, id int, empresaID int) error {
	actual, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if actual.EmpresaID != empresaID {
		return fmt.Errorf("%w: no autorizado para eliminar este inmueble", domain.ErrForbidden)
	}
	for _, unidad := range actual.Unidades {
		if unidad.Estado == "ocupado" || unidad.Estado == "reservado" {
			return fmt.Errorf("no se puede eliminar el inmueble porque tiene unidades ocupadas o reservadas")
		}
	}
	return s.repo.Eliminar(ctx, id)
}

func (s *InmuebleService) ListarUnidades(ctx context.Context, propiedadID int, empresaID int) ([]*domain.Unidad, error) {
	if _, err := s.Obtener(ctx, propiedadID, empresaID); err != nil {
		return nil, err
	}
	return s.repo.ListarUnidades(ctx, propiedadID)
}

func (s *InmuebleService) ObtenerUnidad(ctx context.Context, propiedadID int, unidadID int, empresaID int) (*domain.Unidad, error) {
	if _, err := s.Obtener(ctx, propiedadID, empresaID); err != nil {
		return nil, err
	}
	unidad, err := s.repo.BuscarUnidadPorID(ctx, unidadID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if unidad.PropiedadID != propiedadID {
		return nil, fmt.Errorf("%w: la unidad no pertenece al inmueble indicado", domain.ErrForbidden)
	}
	return unidad, nil
}

func (s *InmuebleService) CrearUnidad(ctx context.Context, propiedadID int, empresaID int, unidad *domain.Unidad) (*domain.Unidad, error) {
	if _, err := s.Obtener(ctx, propiedadID, empresaID); err != nil {
		return nil, err
	}
	unidad.PropiedadID = propiedadID
	if unidad.Tipo == "" {
		unidad.Tipo = "cuarto"
	}
	if unidad.Moneda == "" {
		unidad.Moneda = "PEN"
	}
	if unidad.Estado == "" {
		unidad.Estado = "disponible"
	}
	if unidad.Capacidad <= 0 {
		unidad.Capacidad = 1
	}
	return s.repo.CrearUnidad(ctx, unidad)
}

func (s *InmuebleService) ActualizarUnidad(ctx context.Context, propiedadID int, empresaID int, unidad *domain.Unidad) (*domain.Unidad, error) {
	if _, err := s.Obtener(ctx, propiedadID, empresaID); err != nil {
		return nil, err
	}
	actual, err := s.repo.BuscarUnidadPorID(ctx, unidad.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if actual.PropiedadID != propiedadID {
		return nil, fmt.Errorf("%w: no autorizado para actualizar esta unidad", domain.ErrForbidden)
	}
	unidad.PropiedadID = propiedadID
	if unidad.Tipo == "" {
		unidad.Tipo = actual.Tipo
	}
	if unidad.Moneda == "" {
		unidad.Moneda = actual.Moneda
	}
	if unidad.Estado == "" {
		unidad.Estado = actual.Estado
	}
	if unidad.Capacidad <= 0 {
		unidad.Capacidad = actual.Capacidad
	}
	return s.repo.ActualizarUnidad(ctx, unidad)
}

func (s *InmuebleService) EliminarUnidad(ctx context.Context, propiedadID int, unidadID int, empresaID int) error {
	if _, err := s.Obtener(ctx, propiedadID, empresaID); err != nil {
		return err
	}
	unidad, err := s.repo.BuscarUnidadPorID(ctx, unidadID)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if unidad.PropiedadID != propiedadID {
		return fmt.Errorf("%w: no autorizado para eliminar esta unidad", domain.ErrForbidden)
	}
	return s.repo.EliminarUnidad(ctx, unidadID)
}
