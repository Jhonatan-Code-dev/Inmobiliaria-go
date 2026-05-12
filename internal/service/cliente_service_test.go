package service

import (
	"context"
	"testing"

	"rentals-go/internal/domain"
)

type clienteRepoStub struct {
	cliente *domain.Cliente
	created *domain.Cliente
}

func (s *clienteRepoStub) ListarPaginado(ctx context.Context, filtros domain.ClienteFiltros) ([]*domain.Cliente, int, error) {
	if s.cliente != nil && s.cliente.DocumentoNumero == filtros.Busqueda {
		return []*domain.Cliente{s.cliente}, 1, nil
	}
	return []*domain.Cliente{}, 0, nil
}

func (s *clienteRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Cliente, error) {
	if s.cliente != nil && s.cliente.ID == id {
		return s.cliente, nil
	}
	return nil, domain.ErrNotFound
}

func (s *clienteRepoStub) Crear(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	s.created = c
	return &domain.Cliente{
		ID:                   15,
		EmpresaID:            c.EmpresaID,
		TipoIdentificacionID: c.TipoIdentificacionID,
		DocumentoNumero:      c.DocumentoNumero,
		Nombres:              c.Nombres,
		Estado:               c.Estado,
	}, nil
}

func (s *clienteRepoStub) Actualizar(ctx context.Context, c *domain.Cliente) (*domain.Cliente, error) {
	return c, nil
}

func (s *clienteRepoStub) Eliminar(ctx context.Context, id int) error {
	return nil
}

type tipoIdentificacionRepoStub struct {
	exists bool
	list   []*domain.TipoIdentificacion
}

func (s *tipoIdentificacionRepoStub) ListarActivos(ctx context.Context) ([]*domain.TipoIdentificacion, error) {
	return s.list, nil
}

func (s *tipoIdentificacionRepoStub) ExisteActivo(ctx context.Context, id int) (bool, error) {
	return s.exists, nil
}

func TestRegistrarClienteValidaTipoIdentificacion(t *testing.T) {
	t.Parallel()

	repo := &clienteRepoStub{}
	svc := NewClienteService(repo, &tipoIdentificacionRepoStub{exists: false})

	_, err := svc.RegistrarCliente(context.Background(), &domain.Cliente{
		EmpresaID:            1,
		TipoIdentificacionID: 999,
		DocumentoNumero:      "77889966",
		Nombres:              "Juan",
		Estado:               "activo",
	})
	if err == nil {
		t.Fatal("expected validation error for invalid tipo_identificacion_id")
	}
	if repo.created != nil {
		t.Fatal("expected cliente not to be persisted")
	}
}

func TestRegistrarClienteCreaCuandoTipoIdentificacionExiste(t *testing.T) {
	t.Parallel()

	repo := &clienteRepoStub{}
	svc := NewClienteService(repo, &tipoIdentificacionRepoStub{exists: true})

	out, err := svc.RegistrarCliente(context.Background(), &domain.Cliente{
		EmpresaID:            1,
		TipoIdentificacionID: 1,
		DocumentoNumero:      "77889966",
		Nombres:              "Juan",
		Estado:               "activo",
	})
	if err != nil {
		t.Fatalf("RegistrarCliente() error = %v", err)
	}
	if out.ID != 15 {
		t.Fatalf("returned ID = %d, want 15", out.ID)
	}
	if repo.created == nil {
		t.Fatal("expected cliente to be persisted")
	}
}
