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
	perfilUser, emp, err := s.usuarioRepo.BuscarPerfil(ctx, u.ID)
	if err != nil {
		return "", nil, nil, err
	}
	token, err := auth.GenerarToken(s.jwtSecret, 24*time.Hour, auth.Claims{
		UsuarioID: perfilUser.ID,
		EmpresaID: perfilUser.EmpresaID,
		Rol:       "usuario",
	})
	if err != nil {
		return "", nil, nil, err
	}
	return token, perfilUser, emp, nil
}

func (s *UsuarioService) Perfil(ctx context.Context, id int) (*domain.Usuario, *domain.Empresa, error) {
	return s.usuarioRepo.BuscarPerfil(ctx, id)
}

func (s *UsuarioService) ActualizarPassword(ctx context.Context, id int, password string) error {
	hasher := security.NewServicioHash()
	hash, err := hasher.Encriptar(password)
	if err != nil {
		return err
	}
	return s.usuarioRepo.ActualizarPassword(ctx, id, hash)
}
