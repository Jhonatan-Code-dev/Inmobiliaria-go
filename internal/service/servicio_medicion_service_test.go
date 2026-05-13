package service

import (
	"context"
	"rentals-go/internal/domain"
	"testing"
	"time"
)

type medicionRepoStub struct {
	ultima   *domain.ServicioMedicion
	created  *domain.ServicioMedicion
	updated  *domain.ServicioMedicion
}

func (s *medicionRepoStub) Listar(ctx context.Context, filtros domain.ServicioMedicionFiltros) ([]*domain.ServicioMedicion, int, error) {
	return nil, 0, nil
}
func (s *medicionRepoStub) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.ServicioMedicion, error) {
	if s.created != nil && s.created.ID == id {
		return s.created, nil
	}
	return nil, nil
}
func (s *medicionRepoStub) Crear(ctx context.Context, m *domain.ServicioMedicion) (*domain.ServicioMedicion, error) {
	m.ID = 100
	s.created = m
	return m, nil
}
func (s *medicionRepoStub) Actualizar(ctx context.Context, m *domain.ServicioMedicion) (*domain.ServicioMedicion, error) {
	s.updated = m
	return m, nil
}
func (s *medicionRepoStub) Eliminar(ctx context.Context, id int, empresaID int) error {
	return nil
}
func (s *medicionRepoStub) ObtenerUltimaLectura(ctx context.Context, contratoID int, tipo string) (*domain.ServicioMedicion, error) {
	return s.ultima, nil
}
func (s *medicionRepoStub) BuscarPorFecha(ctx context.Context, unidadID int, tipo string, fecha time.Time) (*domain.ServicioMedicion, error) {
	return nil, nil // Por defecto no hay duplicados en los tests existentes
}

type cargoRepoStub struct {
	created *domain.Cargo
}

func (s *cargoRepoStub) Listar(ctx context.Context, filtros domain.CargoFiltros) ([]*domain.Cargo, int, error) {
	return nil, 0, nil
}
func (s *cargoRepoStub) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Cargo, error) {
	return nil, nil
}
func (s *cargoRepoStub) Crear(ctx context.Context, c *domain.Cargo) (*domain.Cargo, error) {
	c.ID = 500
	s.created = c
	return c, nil
}
func (s *cargoRepoStub) Actualizar(ctx context.Context, c *domain.Cargo) (*domain.Cargo, error) {
	return nil, nil
}
func (s *cargoRepoStub) Eliminar(ctx context.Context, id int, empresaID int) error {
	return nil
}

func TestRegistrarYCobrar(t *testing.T) {
	medRepo := &medicionRepoStub{
		ultima: &domain.ServicioMedicion{LecturaActual: 100},
	}
	carRepo := &cargoRepoStub{}
	alqRepo := &alquilerRepoStub{
		alq: &domain.Alquiler{ID: 1, Moneda: "PEN", UnidadID: 10},
	}

	svc := NewServicioMedicionService(medRepo, carRepo, alqRepo)

	reg := &domain.RegistroLectura{
		ContratoID:     1,
		TipoServicio:   "luz",
		LecturaActual:  150,
		PrecioUnitario: 2.0,
		Factor:         10.0,
		CargoFijo:      5.0,
		FechaLectura:   "2026-05-12",
	}

	med, err := svc.RegistrarYCobrar(context.Background(), reg, 1)
	if err != nil {
		t.Fatalf("RegistrarYCobrar() error = %v", err)
	}

	// Verificar medición
	// (150 - 100) * 10 = 500 unidades
	if med.Consumo != 500 {
		t.Errorf("Consumo = %.2f, want 500", med.Consumo)
	}
	// (500 * 2.0) + 5.0 = 1005.00
	if med.Monto != 1005 {
		t.Errorf("Monto = %.2f, want 1005", med.Monto)
	}
	if !med.Procesado {
		t.Error("expected medicion to be marked as processed")
	}
	if med.CargoID == nil || *med.CargoID != 500 {
		t.Errorf("CargoID = %v, want 500", med.CargoID)
	}

	// Verificar cargo generado
	if carRepo.created == nil {
		t.Fatal("expected cargo to be created")
	}
	if carRepo.created.Monto != 1005 {
		t.Errorf("Cargo Monto = %.2f, want 1005", carRepo.created.Monto)
	}
	if carRepo.created.Moneda != "PEN" {
		t.Errorf("Cargo Moneda = %s, want PEN", carRepo.created.Moneda)
	}
}
