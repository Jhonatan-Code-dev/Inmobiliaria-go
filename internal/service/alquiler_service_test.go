package service

import (
	"context"
	"testing"
	"time"

	"rentals-go/internal/domain"
)

type alquilerRepoStub struct {
	created *domain.Alquiler
}

func (s *alquilerRepoStub) ListarPaginado(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	return nil, 0, nil
}

func (s *alquilerRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Alquiler, error) {
	return nil, domain.ErrNotFound
}

func (s *alquilerRepoStub) Crear(ctx context.Context, alquiler *domain.Alquiler) (*domain.Alquiler, error) {
	s.created = alquiler
	return &domain.Alquiler{ID: 11, EmpresaID: alquiler.EmpresaID, Estado: alquiler.Estado, Tipo: alquiler.Tipo, Moneda: alquiler.Moneda}, nil
}

func (s *alquilerRepoStub) Actualizar(ctx context.Context, alq *domain.Alquiler) (*domain.Alquiler, error) {
	return nil, nil
}

func (s *alquilerRepoStub) Eliminar(ctx context.Context, id int) error {
	return nil
}

type pagoRepoStub struct{}

func (s *pagoRepoStub) Registrar(ctx context.Context, pago *domain.RegistroPagoAlquiler) (*domain.PagoAlquiler, error) {
	return &domain.PagoAlquiler{ID: 9, ContratoID: pago.ContratoID}, nil
}

func (s *pagoRepoStub) ListarPendientesMesActual(ctx context.Context, empresaID int, now time.Time) ([]*domain.PagoPendiente, error) {
	return nil, nil
}

func (s *pagoRepoStub) Listar(ctx context.Context, filtros domain.PagoFiltros) ([]*domain.PagoAlquiler, int, error) {
	return nil, 0, nil
}

func (s *pagoRepoStub) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.PagoAlquiler, error) {
	return nil, nil
}

func (s *pagoRepoStub) Eliminar(ctx context.Context, id int, empresaID int) error {
	return nil
}

func (s *pagoRepoStub) Actualizar(ctx context.Context, pago *domain.PagoAlquiler) (*domain.PagoAlquiler, error) {
	return nil, nil
}

func TestCrearAlquilerAsignaDefaults(t *testing.T) {
	t.Parallel()

	repo := &alquilerRepoStub{}
	svc := NewAlquilerService(repo)

	out, err := svc.Crear(context.Background(), &domain.Alquiler{
		EmpresaID: 1,
		ClienteID: 2,
		UnidadID:  3,
	})
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}
	if out.ID != 11 {
		t.Fatalf("returned ID = %d, want 11", out.ID)
	}
	if repo.created == nil {
		t.Fatal("expected alquiler to be created")
	}
	if repo.created.Moneda != "PEN" || repo.created.Tipo != "alquiler" || repo.created.Estado != "activo" {
		t.Fatalf("defaults not applied: moneda=%q tipo=%q estado=%q", repo.created.Moneda, repo.created.Tipo, repo.created.Estado)
	}
}

func TestRegistrarPagoValidaMesCorrespondiente(t *testing.T) {
	t.Parallel()

	svc := NewPagoAlquilerService(&pagoRepoStub{})
	_, err := svc.Registrar(context.Background(), &domain.RegistroPagoAlquiler{
		EmpresaID:          1,
		ContratoID:         2,
		FechaPago:          time.Now(),
		MontoPagadoCents:   1000,
		MesCorrespondiente: 13,
	})
	if err == nil {
		t.Fatal("expected validation error for invalid mes_correspondiente")
	}
}
