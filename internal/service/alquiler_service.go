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

func (s *AlquilerService) Actualizar(ctx context.Context, id int, empresaID int, alq *domain.Alquiler) (*domain.Alquiler, error) {
	old, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	alq.ID = id
	alq.EmpresaID = empresaID
	// Mantener el código original
	alq.Codigo = old.Codigo
	return s.repo.Actualizar(ctx, alq)
}

func (s *AlquilerService) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return s.repo.Eliminar(ctx, id)
}

func (s *AlquilerService) Terminar(ctx context.Context, id int, empresaID int) error {
	alq, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	alq.Estado = "finalizado"
	alq.ActivoParaCobro = false
	// Al actualizar a finalizado, el repo también debería liberar la unidad o podemos hacerlo explícito.
	// En mi implementación de repo.Eliminar lo hace, pero aquí es una actualización de estado.
	// Ajustaré el repo.Actualizar o lo haré aquí.
	_, err = s.repo.Actualizar(ctx, alq)
	return err
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

func (s *PagoAlquilerService) ListarHistorial(ctx context.Context, empresaID int, pagina, limite int) ([]*domain.PagoAlquiler, int, error) {
	if pagina <= 0 {
		pagina = 1
	}
	if limite <= 0 {
		limite = 10
	}
	return s.repo.Listar(ctx, empresaID, pagina, limite)
}

func (s *PagoAlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.PagoAlquiler, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Anular(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Actualizar(ctx context.Context, id int, empresaID int, notas *string, metodoPago string) (*domain.PagoAlquiler, error) {
	pago, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	if notas != nil {
		pago.Nota = notas
	}
	if metodoPago != "" {
		pago.MetodoPago = metodoPago
	}
	return s.repo.Actualizar(ctx, pago)
}
