package service

import (
	"context"
	"time"

	"rentals-go/internal/domain"
)

// DashboardService implementa domain.DashboardService.
type DashboardService struct {
	repo domain.DashboardRepository
}

func NewDashboardService(repo domain.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

// ObtenerResumenGeneral devuelve los KPIs globales del negocio para el mes actual.
func (s *DashboardService) ObtenerResumenGeneral(ctx context.Context, empresaID int) (*domain.ResumenGeneral, error) {
	return s.repo.ResumenGeneral(ctx, empresaID, time.Now().UTC())
}

// ObtenerOcupacion devuelve la tasa de ocupación global y desglosada por propiedad.
func (s *DashboardService) ObtenerOcupacion(ctx context.Context, empresaID int) (*domain.ResumenOcupacion, error) {
	return s.repo.ResumenOcupacion(ctx, empresaID)
}

// ObtenerMorosidad devuelve la lista de inquilinos con cargos vencidos sin pagar.
func (s *DashboardService) ObtenerMorosidad(ctx context.Context, empresaID int) (*domain.ResumenMorosidad, error) {
	return s.repo.ResumenMorosidad(ctx, empresaID, time.Now().UTC())
}

// ObtenerReporteFinanciero devuelve el balance de ingresos y gastos en un rango de fechas.
// Si desde o hasta están vacíos, usa los últimos 6 meses por defecto.
func (s *DashboardService) ObtenerReporteFinanciero(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ReporteFinanciero, error) {
	ahora := time.Now().UTC()
	if desde.IsZero() {
		desde = time.Date(ahora.Year(), ahora.Month()-5, 1, 0, 0, 0, 0, time.UTC)
	}
	if hasta.IsZero() {
		hasta = ahora
	}
	// Asegurar que desde <= hasta
	if desde.After(hasta) {
		desde, hasta = hasta, desde
	}
	return s.repo.ReporteFinanciero(ctx, empresaID, desde, hasta)
}

// ObtenerContratosProximosVencer devuelve contratos que vencen en los próximos N días.
// Si dias <= 0 usa 30 días como defecto.
func (s *DashboardService) ObtenerContratosProximosVencer(ctx context.Context, empresaID int, dias int) ([]domain.ContratoProximoVencer, error) {
	if dias <= 0 {
		dias = 30
	}
	return s.repo.ContratosProximosVencer(ctx, empresaID, dias, time.Now().UTC())
}

// ObtenerEstadoCuentaCliente devuelve el estado de cuenta detallado de un inquilino.
func (s *DashboardService) ObtenerEstadoCuentaCliente(ctx context.Context, empresaID, clienteID int) (*domain.EstadoCuentaCliente, error) {
	return s.repo.EstadoCuentaCliente(ctx, empresaID, clienteID)
}

// ObtenerTopUnidades devuelve las unidades que más ingresos generaron en un periodo.
// Si desde o hasta están vacíos usa el mes actual. Si limite <= 0 usa 10.
func (s *DashboardService) ObtenerTopUnidades(ctx context.Context, empresaID int, desde, hasta time.Time, limite int) ([]domain.TopUnidad, error) {
	ahora := time.Now().UTC()
	if desde.IsZero() {
		desde = time.Date(ahora.Year(), ahora.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	if hasta.IsZero() {
		hasta = ahora
	}
	if limite <= 0 {
		limite = 10
	}
	return s.repo.TopUnidades(ctx, empresaID, desde, hasta, limite)
}
