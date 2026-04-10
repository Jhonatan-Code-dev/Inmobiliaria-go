package service

import (
	"context"
	"testing"

	"rentals-go/config/security"
	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/auth"
)

type usuarioRepoLoginStub struct {
	hash string
}

func (s *usuarioRepoLoginStub) Crear(ctx context.Context, u *domain.Usuario) (*domain.Usuario, error) {
	return nil, nil
}

func (s *usuarioRepoLoginStub) BuscarPorUsuario(ctx context.Context, usuario string) (*domain.Usuario, error) {
	return &domain.Usuario{
		ID:             7,
		Usuario:        usuario,
		HashContrasena: s.hash,
		Estado:         true,
	}, nil
}

func (s *usuarioRepoLoginStub) BuscarPerfil(ctx context.Context, id int) (*domain.Usuario, *domain.Empresa, error) {
	return &domain.Usuario{
			ID:        id,
			Usuario:   "yonas",
			Estado:    true,
			EmpresaID: 25,
		}, &domain.Empresa{
			ID:     25,
			Nombre: "Empresa Demo",
		}, nil
}

func (s *usuarioRepoLoginStub) ActualizarPassword(ctx context.Context, id int, hashContrasena string) error {
	return nil
}

func TestLoginUsaEmpresaIDDelPerfilEnElToken(t *testing.T) {
	t.Parallel()

	hash, err := security.NewServicioHash().Encriptar("123456")
	if err != nil {
		t.Fatalf("Encriptar() error = %v", err)
	}

	svc := NewUsuarioService(&usuarioRepoLoginStub{hash: hash}, []byte("test-secret"))

	token, user, _, err := svc.Login(context.Background(), "yonas", "123456")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if user.EmpresaID != 25 {
		t.Fatalf("user.EmpresaID = %d, want 25", user.EmpresaID)
	}

	claims, err := auth.ValidarToken(token, []byte("test-secret"))
	if err != nil {
		t.Fatalf("ValidarToken() error = %v", err)
	}
	if claims.EmpresaID != 25 {
		t.Fatalf("claims.EmpresaID = %d, want 25", claims.EmpresaID)
	}
}
