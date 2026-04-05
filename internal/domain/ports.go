package domain

import "context"

// Repositorios (puertos) sin dependencias externas.
type (
	AdminRepository interface {
		BuscarPorUsuario(ctx context.Context, usuario string) (*Admin, error)
		BuscarPorID(ctx context.Context, id int) (*Admin, error)
		ActualizarCredenciales(ctx context.Context, id int, usuario, hashContrasena string) (*Admin, error)
	}

	EmpresaRepository interface {
		Listar(ctx context.Context) ([]*Empresa, error)
		BuscarPorID(ctx context.Context, id int) (*Empresa, error)
		Crear(ctx context.Context, emp *Empresa) (*Empresa, error)
		Actualizar(ctx context.Context, emp *Empresa) (*Empresa, error)
		Eliminar(ctx context.Context, id int) error
	}

	RolRepository interface {
		BuscarPorNombre(ctx context.Context, nombre string) (*Rol, error)
		Crear(ctx context.Context, rol *Rol) (*Rol, error)
	}

	UsuarioRepository interface {
		Crear(ctx context.Context, u *Usuario) (*Usuario, error)
		BuscarPorUsuario(ctx context.Context, usuario string) (*Usuario, error)
		BuscarPerfil(ctx context.Context, id int) (*Usuario, *Empresa, error)
	}

	MembresiaRepository interface {
		AsignarPrincipal(ctx context.Context, empresaID, usuarioID, rolID int) error
	}
)
