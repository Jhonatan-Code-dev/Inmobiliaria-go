package service

import (
	"context"
	"rentals-go/config/security"
	"rentals-go/internal/domain"
)

type StaffService struct {
	repo domain.StaffRepository
}

func NewStaffService(repo domain.StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func (s *StaffService) Listar(ctx context.Context, filtros domain.StaffFiltros) ([]*domain.Staff, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.PorPagina <= 0 {
		filtros.PorPagina = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *StaffService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Staff, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *StaffService) Registrar(ctx context.Context, reg *domain.RegistroStaff) (*domain.Staff, error) {
	hasher := security.NewServicioHash()
	hash, err := hasher.Encriptar(reg.Contrasena)
	if err != nil {
		return nil, err
	}
	return s.repo.Crear(ctx, reg, hash)
}

func (s *StaffService) Actualizar(ctx context.Context, id int, empresaID int, rolID int, estado string) (*domain.Staff, error) {
	return s.repo.Actualizar(ctx, id, empresaID, rolID, estado)
}

func (s *StaffService) Eliminar(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}
