package service

import (
	"context"
	"testing"
	"time"

	"rentals-go/internal/domain"
)

type alquilerRepoStub struct {
	alq              *domain.Alquiler
	created          *domain.Alquiler
	plantilla        *domain.PlantillaContrato
	createdPlantilla *domain.PlantillaContrato
}

func (s *alquilerRepoStub) ListarPaginado(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	return nil, 0, nil
}

func (s *alquilerRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Alquiler, error) {
	if s.alq != nil && s.alq.ID == id {
		return s.alq, nil
	}
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


func (s *alquilerRepoStub) ListarPlantillas(ctx context.Context, empresaID int) ([]*domain.PlantillaContrato, error) {
	return nil, nil
}
func (s *alquilerRepoStub) ObtenerPlantilla(ctx context.Context, id int, empresaID int) (*domain.PlantillaContrato, error) {
	if s.plantilla != nil && s.plantilla.ID == id {
		return s.plantilla, nil
	}
	return nil, domain.ErrNotFound
}
func (s *alquilerRepoStub) CrearPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	return nil, nil
}
func (s *alquilerRepoStub) ActualizarPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	return nil, nil
}
func (s *alquilerRepoStub) EliminarPlantilla(ctx context.Context, id int, empresaID int) error {
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
	svc := NewAlquilerService(repo, &clienteRepoStub{}, &empresaRepoStub{})

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

func TestGenerarContrato(t *testing.T) {
	t.Parallel()

	apellidos := "Pérez"
	clienteMock := &domain.Cliente{
		ID:              2,
		Nombres:         "Juan",
		Apellidos:       &apellidos,
		DocumentoNumero: "12345678",
	}

	alqMock := &domain.Alquiler{
		ID:             1,
		EmpresaID:      1,
		ClienteID:      2,
		UnidadCodigo:   "A-101",
		MontoRenta:     1500.00,
		MontoDeposito:  3000.00,
		Moneda:         "PEN",
		FechaInicio:    time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		DiaVencimiento: 5,
		ClienteNombre:  "Juan",
	}

	plantillaMock := &domain.PlantillaContrato{
		ID:        99,
		EmpresaID: 1,
		Contenido: "Contrato de {{cliente_nombre}} {{cliente_apellidos}} para unidad {{unidad_codigo}} con deposito de {{monto_deposito}}",
	}

	repo := &alquilerRepoStub{alq: alqMock, plantilla: plantillaMock}
	cliRepo := &clienteRepoStub{cliente: clienteMock}
	svc := NewAlquilerService(repo, cliRepo, &empresaRepoStub{})

	// Probar generación con plantilla por defecto
	docDefault, err := svc.GenerarContrato(context.Background(), 1, 1, 0)
	if err != nil {
		t.Fatalf("GenerarContrato() default error = %v", err)
	}
	if len(docDefault) == 0 {
		t.Fatal("expected default content to be generated")
	}

	// Probar generación con plantilla custom
	docCustom, err := svc.GenerarContrato(context.Background(), 1, 1, 99)
	if err != nil {
		t.Fatalf("GenerarContrato() custom error = %v", err)
	}
	expected := "Contrato de Juan Pérez para unidad A-101 con deposito de 3000.00"
	if docCustom != expected {
		t.Errorf("GenerarContrato() = %q, want %q", docCustom, expected)
	}
}

func TestGenerarContratoBorrador(t *testing.T) {
	t.Parallel()

	apellidos := "DB Apellidos"
	clienteMock := &domain.Cliente{
		ID:              2,
		Nombres:         "DB Nombre",
		Apellidos:       &apellidos,
		DocumentoNumero: "12345678",
	}

	plantillaMock := &domain.PlantillaContrato{
		ID:        99,
		EmpresaID: 1,
		Contenido: "Borrador de {{cliente_nombre}} {{cliente_apellidos}} con DNI {{cliente_documento}} y renta {{monto_renta}}",
	}

	repo := &alquilerRepoStub{plantilla: plantillaMock}
	cliRepo := &clienteRepoStub{cliente: clienteMock}
	svc := NewAlquilerService(repo, cliRepo, &empresaRepoStub{})

	req := domain.GenerarBorradorRequest{
		PlantillaID:      99,
		ClienteDocumento: "12345678",
		ClienteNombre:    "Fallback Name", // Debería ser sobrescrito por DB ("DB Nombre")
		UnidadCodigo:     "A-101",
		MontoRenta:       1500,
	}

	doc, err := svc.GenerarContratoBorrador(context.Background(), 1, req)
	if err != nil {
		t.Fatalf("GenerarContratoBorrador() error = %v", err)
	}

	if len(doc) == 0 {
		t.Errorf("GenerarContratoBorrador() returned empty document")
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
