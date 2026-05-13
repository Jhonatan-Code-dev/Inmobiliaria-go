package service

import (
	"context"
	"rentals-go/internal/domain"
	"testing"
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
		FechaLectura:   "2026-05-12",
	}

	med, err := svc.RegistrarYCobrar(context.Background(), reg, 1)
	if err != nil {
		t.Fatalf("RegistrarYCobrar() error = %v", err)
	}

	// Verificar medición
	if med.Consumo != 50 {
		t.Errorf("Consumo = %.2f, want 50", med.Consumo)
	}
	if med.Monto != 100 {
		t.Errorf("Monto = %.2f, want 100", med.Monto)
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
	if carRepo.created.Monto != 100 {
		t.Errorf("Cargo Monto = %.2f, want 100", carRepo.created.Monto)
	}
	if carRepo.created.Moneda != "PEN" {
		t.Errorf("Cargo Moneda = %s, want PEN", carRepo.created.Moneda)
	}
}
