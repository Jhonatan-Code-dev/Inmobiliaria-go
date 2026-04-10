package service

import (
	"context"
	"testing"

	"rentals-go/internal/domain"
)

type inmuebleRepoStub struct {
	inmueble      *domain.Inmueble
	unidad        *domain.Unidad
	created       *domain.Inmueble
	createdUnidad *domain.Unidad
}

func (s *inmuebleRepoStub) ListarPaginado(ctx context.Context, filtros domain.InmuebleFiltros) ([]*domain.Inmueble, int, error) {
	return []*domain.Inmueble{}, 0, nil
}

func (s *inmuebleRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Inmueble, error) {
	if s.inmueble != nil && s.inmueble.ID == id {
		return s.inmueble, nil
	}
	return nil, domain.ErrNotFound
}

func (s *inmuebleRepoStub) Crear(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	s.created = inmueble
	return &domain.Inmueble{
		ID:            21,
		EmpresaID:     inmueble.EmpresaID,
		Nombre:        inmueble.Nombre,
		Tipo:          inmueble.Tipo,
		Estado:        inmueble.Estado,
		TotalPisos:    inmueble.TotalPisos,
		TotalUnidades: inmueble.TotalUnidades,
	}, nil
}

func (s *inmuebleRepoStub) Actualizar(ctx context.Context, inmueble *domain.Inmueble) (*domain.Inmueble, error) {
	return inmueble, nil
}

func (s *inmuebleRepoStub) Eliminar(ctx context.Context, id int) error { return nil }

func (s *inmuebleRepoStub) ListarUnidades(ctx context.Context, propiedadID int) ([]*domain.Unidad, error) {
	return []*domain.Unidad{}, nil
}

func (s *inmuebleRepoStub) BuscarUnidadPorID(ctx context.Context, id int) (*domain.Unidad, error) {
	if s.unidad != nil && s.unidad.ID == id {
		return s.unidad, nil
	}
	return nil, domain.ErrNotFound
}

func (s *inmuebleRepoStub) CrearUnidad(ctx context.Context, unidad *domain.Unidad) (*domain.Unidad, error) {
	s.createdUnidad = unidad
	return &domain.Unidad{
		ID:                34,
		PropiedadID:       unidad.PropiedadID,
		Codigo:            unidad.Codigo,
		Tipo:              unidad.Tipo,
		Moneda:            unidad.Moneda,
		Estado:            unidad.Estado,
		Capacidad:         unidad.Capacidad,
		DepositoRequerido: unidad.DepositoRequerido,
	}, nil
}

func (s *inmuebleRepoStub) ActualizarUnidad(ctx context.Context, unidad *domain.Unidad) (*domain.Unidad, error) {
	return unidad, nil
}

func (s *inmuebleRepoStub) EliminarUnidad(ctx context.Context, id int) error { return nil }

func TestCrearInmuebleAsignaDefaults(t *testing.T) {
	t.Parallel()

	repo := &inmuebleRepoStub{}
	svc := NewInmuebleService(repo)

	out, err := svc.Crear(context.Background(), &domain.Inmueble{
		EmpresaID: 1,
		Nombre:    "Edificio Central",
		Direccion: "Av. Principal 123",
	})
	if err != nil {
		t.Fatalf("Crear() error = %v", err)
	}
	if out.ID != 21 {
		t.Fatalf("returned ID = %d, want 21", out.ID)
	}
	if repo.created == nil {
		t.Fatal("expected inmueble to be persisted")
	}
	if repo.created.Tipo != "casa" || repo.created.Estado != "activa" {
		t.Fatalf("defaults not applied: tipo=%q estado=%q", repo.created.Tipo, repo.created.Estado)
	}
}

func TestCrearUnidadAsignaDefaults(t *testing.T) {
	t.Parallel()

	repo := &inmuebleRepoStub{
		inmueble: &domain.Inmueble{ID: 8, EmpresaID: 1, Nombre: "Casa 1"},
	}
	svc := NewInmuebleService(repo)

	out, err := svc.CrearUnidad(context.Background(), 8, 1, &domain.Unidad{
		Codigo: "A-101",
	})
	if err != nil {
		t.Fatalf("CrearUnidad() error = %v", err)
	}
	if out.ID != 34 {
		t.Fatalf("returned ID = %d, want 34", out.ID)
	}
	if repo.createdUnidad == nil {
		t.Fatal("expected unidad to be persisted")
	}
	if repo.createdUnidad.Moneda != "PEN" || repo.createdUnidad.Estado != "disponible" || repo.createdUnidad.Tipo != "cuarto" {
		t.Fatalf("defaults not applied: moneda=%q estado=%q tipo=%q", repo.createdUnidad.Moneda, repo.createdUnidad.Estado, repo.createdUnidad.Tipo)
	}
}

func TestObtenerUnidadValidaPertenencia(t *testing.T) {
	t.Parallel()

	repo := &inmuebleRepoStub{
		inmueble: &domain.Inmueble{ID: 8, EmpresaID: 1},
		unidad:   &domain.Unidad{ID: 12, PropiedadID: 9},
	}
	svc := NewInmuebleService(repo)

	_, err := svc.ObtenerUnidad(context.Background(), 8, 12, 1)
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}
