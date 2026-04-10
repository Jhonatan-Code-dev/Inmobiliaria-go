package service

import (
	"context"
	"fmt"
	"time"

	"rentals-go/internal/domain"
)

type AlquilerService struct {
	repo domain.AlquilerRepository
}

func NewAlquilerService(repo domain.AlquilerRepository) *AlquilerService {
	return &AlquilerService{repo: repo}
}

func (s *AlquilerService) Listar(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *AlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Alquiler, error) {
	item, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if item.EmpresaID != empresaID {
		return nil, fmt.Errorf("%w: alquiler no pertenece a la empresa", domain.ErrForbidden)
	}
	return item, nil
}

func (s *AlquilerService) Crear(ctx context.Context, alquiler *domain.Alquiler) (*domain.Alquiler, error) {
	if alquiler.Tipo == "" {
		alquiler.Tipo = "alquiler"
	}
	if alquiler.Moneda == "" {
		alquiler.Moneda = "PEN"
	}
	if alquiler.Estado == "" {
		alquiler.Estado = "activo"
	}
	alquiler.ActivoParaCobro = true
	alquiler.Codigo = fmt.Sprintf("ALQ-%d", time.Now().UTC().UnixNano())
	return s.repo.Crear(ctx, alquiler)
}

type PagoAlquilerService struct {
	repo domain.PagoAlquilerRepository
}

func NewPagoAlquilerService(repo domain.PagoAlquilerRepository) *PagoAlquilerService {
	return &PagoAlquilerService{repo: repo}
}

func (s *PagoAlquilerService) Registrar(ctx context.Context, pago *domain.RegistroPagoAlquiler) (*domain.PagoAlquiler, error) {
	if pago.MesCorrespondiente < 1 || pago.MesCorrespondiente > 12 {
		return nil, fmt.Errorf("mes_correspondiente debe estar entre 1 y 12")
	}
	if pago.FechaPago.IsZero() {
		return nil, fmt.Errorf("fecha_pago es obligatoria")
	}
	return s.repo.Registrar(ctx, pago)
}

func (s *PagoAlquilerService) ListarPendientesMesActual(ctx context.Context, empresaID int) ([]*domain.PagoPendiente, error) {
	return s.repo.ListarPendientesMesActual(ctx, empresaID, time.Now().UTC())
}
