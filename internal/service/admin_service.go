package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"rentals-go/config/security"
	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/auth"
	"rentals-go/internal/pkg/moneda"
)

var (
	ErrCredenciales         = errors.New("credenciales inválidas")
	ErrMonedaInvalida       = errors.New("moneda inválida")
	ErrUsuarioAdminInvalido = errors.New("usuario de admin inválido")
	ErrContrasenaInvalida   = errors.New("contraseña inválida")
)

type AdminService struct {
	adminRepo   domain.AdminRepository
	empresaRepo domain.EmpresaRepository
	usuarioRepo domain.UsuarioRepository
	rolRepo     domain.RolRepository
	membRepo    domain.MembresiaRepository
	jwtSecret   []byte
}

func NewAdminService(adminRepo domain.AdminRepository, empresaRepo domain.EmpresaRepository, usuarioRepo domain.UsuarioRepository, rolRepo domain.RolRepository, membRepo domain.MembresiaRepository, jwtSecret []byte) *AdminService {
	return &AdminService{
		adminRepo:   adminRepo,
		empresaRepo: empresaRepo,
		usuarioRepo: usuarioRepo,
		rolRepo:     rolRepo,
		membRepo:    membRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *AdminService) Login(ctx context.Context, usuario, contrasena string) (string, error) {
	a, err := s.adminRepo.BuscarPorUsuario(ctx, usuario)
	if err != nil || !a.Activo {
		return "", ErrCredenciales
	}
	hasher := security.NewServicioHash()
	if !hasher.Comparar(a.HashContrasena, contrasena) {
		return "", ErrCredenciales
	}
	return auth.GenerarToken(s.jwtSecret, 24*time.Hour, auth.Claims{
		AdminID: a.ID,
		Rol:     "admin",
	})
}

func (s *AdminService) Perfil(ctx context.Context, adminID int) (*domain.Admin, error) {
	return s.adminRepo.BuscarPorID(ctx, adminID)
}

func (s *AdminService) CrearEmpresaConUsuario(ctx context.Context, emp *domain.Empresa, usuario *domain.Usuario, rolID int) (*domain.Empresa, *domain.Usuario, error) {
	emp.Moneda = moneda.NormalizarCodigo(emp.Moneda)
	if err := moneda.ValidarCodigo(emp.Moneda); err != nil {
		return nil, nil, ErrMonedaInvalida
	}
	createdEmp, err := s.empresaRepo.Crear(ctx, emp)
	if err != nil {
		return nil, nil, err
	}
	createdUser, err := s.usuarioRepo.Crear(ctx, usuario)
	if err != nil {
		return nil, nil, err
	}
	roleID, err := s.obtenerRolAdministrador(ctx, rolID)
	if err != nil {
		return nil, nil, err
	}
	if err := s.membRepo.AsignarPrincipal(ctx, createdEmp.ID, createdUser.ID, roleID); err != nil {
		return nil, nil, err
	}
	return createdEmp, createdUser, nil
}

func (s *AdminService) obtenerRolAdministrador(ctx context.Context, rolID int) (int, error) {
	if rolID > 0 {
		return rolID, nil
	}
	rolAdmin, err := s.rolRepo.BuscarPorNombre(ctx, "administrador")
	if err == nil {
		return rolAdmin.ID, nil
	}
	rolNuevo, err := s.rolRepo.Crear(ctx, &domain.Rol{
		Nombre:      "administrador",
		Descripcion: "Administrador de empresa",
	})
	if err != nil {
		return 0, err
	}
	return rolNuevo.ID, nil
}

func (s *AdminService) ListarEmpresas(ctx context.Context) ([]*domain.Empresa, error) {
	return s.empresaRepo.Listar(ctx)
}

func (s *AdminService) ObtenerEmpresa(ctx context.Context, id int) (*domain.Empresa, error) {
	return s.empresaRepo.BuscarPorID(ctx, id)
}

func (s *AdminService) ActualizarEmpresa(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	emp.Moneda = moneda.NormalizarCodigo(emp.Moneda)
	if err := moneda.ValidarCodigo(emp.Moneda); err != nil {
		return nil, ErrMonedaInvalida
	}
	return s.empresaRepo.Actualizar(ctx, emp)
}

func (s *AdminService) EliminarEmpresa(ctx context.Context, id int) error {
	return s.empresaRepo.Eliminar(ctx, id)
}

func (s *AdminService) ActualizarCredenciales(ctx context.Context, adminID int, usuarioNuevo, contrasenaNueva string) (*domain.Admin, error) {
	adminActual, err := s.adminRepo.BuscarPorID(ctx, adminID)
	if err != nil {
		return nil, err
	}

	usuarioNuevo = strings.TrimSpace(usuarioNuevo)
	if usuarioNuevo == "" {
		usuarioNuevo = adminActual.Usuario
	} else if len(usuarioNuevo) < 3 {
		return nil, ErrUsuarioAdminInvalido
	}

	hash := adminActual.HashContrasena
	if contrasenaNueva != "" {
		if len(contrasenaNueva) < 6 {
			return nil, ErrContrasenaInvalida
		}
		hasher := security.NewServicioHash()
		hash, err = hasher.Encriptar(contrasenaNueva)
		if err != nil {
			return nil, err
		}
	}

	return s.adminRepo.ActualizarCredenciales(ctx, adminID, usuarioNuevo, hash)
}
