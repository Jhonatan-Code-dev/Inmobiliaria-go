package service

import (
	"context"
	"time"

	"rentals-go/internal/domain"
)

type ReporteService struct {
	repo domain.ReporteRepository
}

func NewReporteService(repo domain.ReporteRepository) *ReporteService {
	return &ReporteService{repo: repo}
}

// helper to validate and default date ranges to the last 12 months
func (s *ReporteService) validarYCompletarFechas(desde, hasta time.Time) (time.Time, time.Time) {
	ahora := time.Now().UTC()

	// Default to last 12 months if dates are empty
	if desde.IsZero() {
		// Start of month, 11 months ago (total of 12 months including current month)
		desde = time.Date(ahora.Year(), ahora.Month()-11, 1, 0, 0, 0, 0, time.UTC)
	}
	if hasta.IsZero() {
		// End of current month
		hasta = ahora
	}

	// Ensure desde is before or equal to hasta
	if desde.After(hasta) {
		desde, hasta = hasta, desde
	}

	return desde, hasta
}

func (s *ReporteService) ObtenerIngresosGastos(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ReporteIngresosGastos, error) {
	d, h := s.validarYCompletarFechas(desde, hasta)
	return s.repo.ObtenerIngresosGastos(ctx, empresaID, d, h)
}

func (s *ReporteService) ObtenerDistribucionMetodosPago(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionMetodoPago, error) {
	d, h := s.validarYCompletarFechas(desde, hasta)
	return s.repo.ObtenerDistribucionMetodosPago(ctx, empresaID, d, h)
}

func (s *ReporteService) ObtenerDistribucionCategoriasGastos(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionCategoriaGasto, error) {
	d, h := s.validarYCompletarFechas(desde, hasta)
	return s.repo.ObtenerDistribucionCategoriasGastos(ctx, empresaID, d, h)
}

func (s *ReporteService) ObtenerRentabilidadPropiedades(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.RentabilidadPropiedad, error) {
	d, h := s.validarYCompletarFechas(desde, hasta)
	return s.repo.ObtenerRentabilidadPropiedades(ctx, empresaID, d, h)
}

func (s *ReporteService) ObtenerResumenMantenimiento(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ResumenMantenimientoReporte, error) {
	d, h := s.validarYCompletarFechas(desde, hasta)
	return s.repo.ObtenerResumenMantenimiento(ctx, empresaID, d, h)
}
