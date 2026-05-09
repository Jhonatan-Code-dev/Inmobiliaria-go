package service

import (
	"context"
	"testing"
	"time"

	"rentals-go/internal/domain"
)

type gastoRepoStub struct {
	created *domain.Gasto
}

func (s *gastoRepoStub) ListarPaginado(ctx context.Context, filtros domain.GastoFiltros) ([]*domain.Gasto, int, error) {
	return nil, 0, nil
}

func (s *gastoRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Gasto, error) {
	return nil, nil
}

func (s *gastoRepoStub) Crear(ctx context.Context, gasto *domain.Gasto) (*domain.Gasto, error) {
	s.created = gasto
	return &domain.Gasto{
		ID:          9,
		EmpresaID:   gasto.EmpresaID,
		Monto:       gasto.Monto,
		MontoCents:  gasto.MontoCents,
		Fecha:       gasto.Fecha,
		TipoPagoID:  gasto.TipoPagoID,
		Descripcion: gasto.Descripcion,
	}, nil
}

func (s *gastoRepoStub) Actualizar(ctx context.Context, gasto *domain.Gasto) (*domain.Gasto, error) {
	return gasto, nil
}

func (s *gastoRepoStub) Eliminar(ctx context.Context, id int) error {
	return nil
}

type movRepoStub struct {
	created *domain.MovimientoCaja
}

func (s *movRepoStub) Crear(ctx context.Context, mov *domain.MovimientoCaja) (*domain.MovimientoCaja, error) {
	s.created = mov
	return mov, nil
}

type tipoPagoRepoStub struct{}

func (s *tipoPagoRepoStub) Listar(ctx context.Context) ([]*domain.TipoPago, error) {
	return []*domain.TipoPago{{ID: 3, Nombre: "yape"}}, nil
}

func TestRegistrarGastoMantieneCentavosEnMovimiento(t *testing.T) {
	t.Parallel()

	gastoRepo := &gastoRepoStub{}
	movRepo := &movRepoStub{}
	svc := NewGastoService(gastoRepo, movRepo, &tipoPagoRepoStub{})

	fecha := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)
	in := &domain.Gasto{
		EmpresaID:   1,
		Monto:       50.50,
		MontoCents:  5050,
		Fecha:       fecha,
		TipoPagoID:  3,
		Descripcion: "Pago mensual Starlink",
	}

	out, err := svc.RegistrarGasto(context.Background(), in)
	if err != nil {
		t.Fatalf("RegistrarGasto() error = %v", err)
	}
	if out.Monto != 50.50 {
		t.Fatalf("returned monto = %v, want 50.50", out.Monto)
	}
	if movRepo.created == nil {
		t.Fatal("expected movimiento de caja to be created")
	}
	if movRepo.created.Monto != 50.50 {
		t.Fatalf("movimiento monto = %v, want 50.50", movRepo.created.Monto)
	}
}

func TestExportarExcelGeneraBuffer(t *testing.T) {
	t.Parallel()

	gastoRepo := &gastoRepoStubReporte{}
	svc := NewGastoService(gastoRepo, &movRepoStub{}, &tipoPagoRepoStub{})

	filtros := domain.GastoFiltros{EmpresaID: 1}
	buf, err := svc.ExportarExcel(context.Background(), filtros)

	if err != nil {
		t.Fatalf("ExportarExcel() error = %v", err)
	}
	if len(buf) == 0 {
		t.Fatal("expected non-empty buffer for Excel")
	}
}

func TestExportarPDFGeneraBuffer(t *testing.T) {
	t.Parallel()

	gastoRepo := &gastoRepoStubReporte{}
	svc := NewGastoService(gastoRepo, &movRepoStub{}, &tipoPagoRepoStub{})

	filtros := domain.GastoFiltros{EmpresaID: 1}
	buf, err := svc.ExportarPDF(context.Background(), filtros)

	if err != nil {
		t.Fatalf("ExportarPDF() error = %v", err)
	}
	if len(buf) == 0 {
		t.Fatal("expected non-empty buffer for PDF")
	}
}

type gastoRepoStubReporte struct {
	gastoRepoStub
}

func (s *gastoRepoStubReporte) ListarPaginado(ctx context.Context, filtros domain.GastoFiltros) ([]*domain.Gasto, int, error) {
	return []*domain.Gasto{
		{ID: 1, EmpresaID: 1, Monto: 100.0, Fecha: time.Now(), Descripcion: "Gasto Test 1", TipoPagoID: 3},
		{ID: 2, EmpresaID: 1, Monto: 250.50, Fecha: time.Now(), Descripcion: "Gasto Test 2", TipoPagoID: 3},
	}, 2, nil
}
