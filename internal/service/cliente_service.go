package service

import (
	"context"
	"fmt"
	"rentals-go/internal/domain"
)

type ClienteService struct {
	repo      domain.ClienteRepository
	tipoIDRepo domain.TipoIdentificacionRepository
}

func NewClienteService(repo domain.ClienteRepository, tipoIDRepo domain.TipoIdentificacionRepository) *ClienteService {
	return &ClienteService{repo: repo, tipoIDRepo: tipoIDRepo}
}

func (s *ClienteService) Listar(ctx context.Context, filtros domain.ClienteFiltros) ([]*domain.Cliente, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 || filtros.Limite > 10 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *ClienteService) ObtenerCliente(ctx context.Context, id int, empresaID int) (*domain.Cliente, error) {
	c, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if c.EmpresaID != empresaID {
		return nil, fmt.Errorf("%w: cliente no pertenece a la empresa", domain.ErrForbidden)
	}
	return c, nil
}

func (s *ClienteService) ListarTiposIdentificacion(ctx context.Context) ([]*domain.TipoIdentificacion, error) {
	return s.tipoIDRepo.ListarActivos(ctx)
}

func (s *ClienteService) RegistrarCliente(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	ok, err := s.tipoIDRepo.ExisteActivo(ctx, c.TipoIdentificacionID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("tipo_identificacion_id no existe o está inactivo")
	}
	return s.repo.Crear(ctx, c)
}

func (s *ClienteService) ActualizarCliente(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	existente, err := s.repo.BuscarPorID(ctx, c.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if existente.EmpresaID != c.EmpresaID {
		return nil, fmt.Errorf("%w: no autorizado para actualizar este cliente", domain.ErrForbidden)
	}

	ok, err := s.tipoIDRepo.ExisteActivo(ctx, c.TipoIdentificacionID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("tipo_identificacion_id no existe o está inactivo")
	}
	return s.repo.Actualizar(ctx, c)
}

func (s *ClienteService) EliminarCliente(ctx context.Context, id int, empresaID int) error {
	existente, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if existente.EmpresaID != empresaID {
		return fmt.Errorf("%w: no autorizado para eliminar este cliente", domain.ErrForbidden)
	}
	return s.repo.Eliminar(ctx, id)
}
