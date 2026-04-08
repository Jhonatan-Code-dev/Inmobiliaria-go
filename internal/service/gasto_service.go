package service

import (
	"context"
	"fmt"
	"rentals-go/internal/domain"
)

type GastoService struct {
	repo       domain.GastoRepository
	movRepo    domain.MovimientoCajaRepository
	tipoPgRepo domain.TipoPagoRepository
}

func NewGastoService(repo domain.GastoRepository, movRepo domain.MovimientoCajaRepository, tipoPgRepo domain.TipoPagoRepository) *GastoService {
	return &GastoService{
		repo:       repo,
		movRepo:    movRepo,
		tipoPgRepo: tipoPgRepo,
	}
}

func (s *GastoService) ListarTiposPago(ctx context.Context) ([]*domain.TipoPago, error) {
	return s.tipoPgRepo.Listar(ctx)
}

func (s *GastoService) Listar(ctx context.Context, filtros domain.GastoFiltros) ([]*domain.Gasto, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *GastoService) ObtenerGasto(ctx context.Context, id int, empresaID int) (*domain.Gasto, error) {
	g, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if g.EmpresaID != empresaID {
		return nil, fmt.Errorf("gasto no pertenece a la empresa")
	}
	return g, nil
}

func (s *GastoService) RegistrarGasto(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	// 1. Crear el gasto
	nuevoGasto, err := s.repo.Crear(ctx, g)
	if err != nil {
		return nil, err
	}

	// Fetch tipo pago name to use as method for MovimientoCaja
	metodo := "otro" // Default fallback
	tipos, err := s.tipoPgRepo.Listar(ctx)
	if err == nil {
		for _, t := range tipos {
			if t.ID == nuevoGasto.TipoPagoID {
				metodo = t.Nombre
				break
			}
		}
	}

	// 2. Crear el movimiento de caja (egreso)
	mov := &domain.MovimientoCaja{
		EmpresaID:       nuevoGasto.EmpresaID,
		GastoID:         &nuevoGasto.ID,
		Tipo:            "egreso",
		Concepto:        fmt.Sprintf("Gasto: %s", nuevoGasto.Descripcion),
		FechaMovimiento: nuevoGasto.Fecha,
		Moneda:          "PEN", // Por defecto usaremos PEN para el movimiento simplificado de ser necesario
		Monto:           int64(nuevoGasto.Monto),
		Metodo:          metodo,
	}
	_, err = s.movRepo.Crear(ctx, mov)
	if err != nil {
		fmt.Printf("Error al crear movimiento de caja para gasto %d: %v\n", nuevoGasto.ID, err)
	}

	return nuevoGasto, nil
}

func (s *GastoService) ActualizarGasto(ctx context.Context, g *domain.Gasto) (*domain.Gasto, error) {
	// Verificar pertenencia antes de actualizar
	existente, err := s.repo.BuscarPorID(ctx, g.ID)
	if err != nil {
		return nil, err
	}
	if existente.EmpresaID != g.EmpresaID {
		return nil, fmt.Errorf("no autorizado para actualizar este gasto")
	}

	return s.repo.Actualizar(ctx, g)
}

func (s *GastoService) EliminarGasto(ctx context.Context, id int, empresaID int) error {
	existente, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if existente.EmpresaID != empresaID {
		return fmt.Errorf("no autorizado para eliminar este gasto")
	}

	return s.repo.Eliminar(ctx, id)
}
