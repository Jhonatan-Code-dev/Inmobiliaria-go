package service

import (
	"context"
	"time"

	"rentals-go/config/security"
	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/auth"
)

type UsuarioService struct {
	usuarioRepo domain.UsuarioRepository
	jwtSecret   []byte
}

func NewUsuarioService(usuarioRepo domain.UsuarioRepository, jwtSecret []byte) *UsuarioService {
	return &UsuarioService{
		usuarioRepo: usuarioRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *UsuarioService) Login(ctx context.Context, username, contrasena string) (string, *domain.Usuario, *domain.Empresa, error) {
	u, err := s.usuarioRepo.BuscarPorUsuario(ctx, username)
	if err != nil {
		return "", nil, nil, ErrCredenciales
	}
	hasher := security.NewServicioHash()
	if !hasher.Comparar(u.HashContrasena, contrasena) {
		return "", nil, nil, ErrCredenciales
	}
	_, emp, _ := s.usuarioRepo.BuscarPerfil(ctx, u.ID)
	token, err := auth.GenerarToken(s.jwtSecret, 24*time.Hour, auth.Claims{
		UsuarioID: u.ID,
		EmpresaID: u.EmpresaID,
		Rol:       "usuario",
	})
	if err != nil {
		return "", nil, nil, err
	}
	return token, u, emp, nil
}

func (s *UsuarioService) Perfil(ctx context.Context, id int) (*domain.Usuario, *domain.Empresa, error) {
	return s.usuarioRepo.BuscarPerfil(ctx, id)
}
