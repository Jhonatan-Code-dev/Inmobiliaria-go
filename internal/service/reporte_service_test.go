package service

import (
	"context"
	"testing"
	"time"

	"rentals-go/internal/domain"
)

type reporteRepoMock struct {
	llamadoIngresosGastos      bool
	llamadoDistribucionMetodos bool
	llamadoDistribucionGastos  bool
	llamadoRentabilidad        bool
	llamadoMantenimiento       bool
}

func (m *reporteRepoMock) ObtenerIngresosGastos(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ReporteIngresosGastos, error) {
	m.llamadoIngresosGastos = true
	return &domain.ReporteIngresosGastos{
		Desde:         desde.Format("2006-01-02"),
		Hasta:         hasta.Format("2006-01-02"),
		TotalIngresos: 1500.0,
		TotalGastos:   500.0,
		BalanceNeto:   1000.0,
		Serie: []domain.PuntoIngresoGasto{
			{Periodo: "2026-05", Ingresos: 1500.0, Gastos: 500.0, Balance: 1000.0},
		},
	}, nil
}

func (m *reporteRepoMock) ObtenerDistribucionMetodosPago(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionMetodoPago, error) {
	m.llamadoDistribucionMetodos = true
	return []domain.DistribucionMetodoPago{
		{Metodo: "yape", Total: 1000.0, CantidadPagos: 5, Porcentaje: 66.67},
		{Metodo: "efectivo", Total: 500.0, CantidadPagos: 2, Porcentaje: 33.33},
	}, nil
}

func (m *reporteRepoMock) ObtenerDistribucionCategoriasGastos(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.DistribucionCategoriaGasto, error) {
	m.llamadoDistribucionGastos = true
	return []domain.DistribucionCategoriaGasto{
		{TipoPagoID: 1, Categoria: "Servicios", Total: 300.0, CantidadGastos: 2, Porcentaje: 60.0},
		{TipoPagoID: 2, Categoria: "Mantenimiento", Total: 200.0, CantidadGastos: 1, Porcentaje: 40.0},
	}, nil
}

func (m *reporteRepoMock) ObtenerRentabilidadPropiedades(ctx context.Context, empresaID int, desde, hasta time.Time) ([]domain.RentabilidadPropiedad, error) {
	m.llamadoRentabilidad = true
	return []domain.RentabilidadPropiedad{
		{PropiedadID: 1, Nombre: "Condominio A", Direccion: "Av. Sol 123", TotalUnidades: 10, UnidadesOcupadas: 8, TasaOcupacionPct: 80.0, Ingresos: 1000.0, Gastos: 200.0, Rentabilidad: 800.0},
	}, nil
}

func (m *reporteRepoMock) ObtenerResumenMantenimiento(ctx context.Context, empresaID int, desde, hasta time.Time) (*domain.ResumenMantenimientoReporte, error) {
	m.llamadoMantenimiento = true
	return &domain.ResumenMantenimientoReporte{
		TotalTickets: 5,
		PorEstado:    domain.TicketsPorEstado{Abierto: 2, EnProgreso: 2, Resuelto: 1},
		PorPrioridad: domain.TicketsPorPrioridad{Baja: 1, Media: 3, Alta: 1},
	}, nil
}

func TestReporteService_ObtenerIngresosGastos(t *testing.T) {
	mock := &reporteRepoMock{}
	svc := NewReporteService(mock)

	ctx := context.Background()
	res, err := svc.ObtenerIngresosGastos(ctx, 1, time.Time{}, time.Time{})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mock.llamadoIngresosGastos {
		t.Fatal("se esperaba que el repositorio fuera invocado")
	}
	if res.TotalIngresos != 1500.0 || res.TotalGastos != 500.0 || res.BalanceNeto != 1000.0 {
		t.Errorf("totales incorrectos obtenidos: %+v", res)
	}
}

func TestReporteService_ObtenerDistribucionMetodosPago(t *testing.T) {
	mock := &reporteRepoMock{}
	svc := NewReporteService(mock)

	ctx := context.Background()
	res, err := svc.ObtenerDistribucionMetodosPago(ctx, 1, time.Time{}, time.Time{})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mock.llamadoDistribucionMetodos {
		t.Fatal("se esperaba que el repositorio fuera invocado")
	}
	if len(res) != 2 {
		t.Errorf("se esperaban 2 items, se obtuvieron %d", len(res))
	}
}

func TestReporteService_ObtenerDistribucionCategoriasGastos(t *testing.T) {
	mock := &reporteRepoMock{}
	svc := NewReporteService(mock)

	ctx := context.Background()
	res, err := svc.ObtenerDistribucionCategoriasGastos(ctx, 1, time.Time{}, time.Time{})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mock.llamadoDistribucionGastos {
		t.Fatal("se esperaba que el repositorio fuera invocado")
	}
	if len(res) != 2 {
		t.Errorf("se esperaban 2 items, se obtuvieron %d", len(res))
	}
}

func TestReporteService_ObtenerRentabilidadPropiedades(t *testing.T) {
	mock := &reporteRepoMock{}
	svc := NewReporteService(mock)

	ctx := context.Background()
	res, err := svc.ObtenerRentabilidadPropiedades(ctx, 1, time.Time{}, time.Time{})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mock.llamadoRentabilidad {
		t.Fatal("se esperaba que el repositorio fuera invocado")
	}
	if len(res) != 1 || res[0].Nombre != "Condominio A" {
		t.Errorf("propiedad devuelta incorrecta: %+v", res)
	}
}

func TestReporteService_ObtenerResumenMantenimiento(t *testing.T) {
	mock := &reporteRepoMock{}
	svc := NewReporteService(mock)

	ctx := context.Background()
	res, err := svc.ObtenerResumenMantenimiento(ctx, 1, time.Time{}, time.Time{})

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if !mock.llamadoMantenimiento {
		t.Fatal("se esperaba que el repositorio fuera invocado")
	}
	if res.TotalTickets != 5 || res.PorEstado.Abierto != 2 {
		t.Errorf("resumen de mantenimiento incorrecto: %+v", res)
	}
}
