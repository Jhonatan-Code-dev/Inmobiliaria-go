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
		ListarPaginado(ctx context.Context, limite, offset int, busqueda string) ([]*Empresa, int, error)
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

	MovimientoCajaRepository interface {
		Crear(ctx context.Context, mov *MovimientoCaja) (*MovimientoCaja, error)
	}

	GastoRepository interface {
		ListarPaginado(ctx context.Context, filtros GastoFiltros) ([]*Gasto, int, error)
		BuscarPorID(ctx context.Context, id int) (*Gasto, error)
		Crear(ctx context.Context, gasto *Gasto) (*Gasto, error)
		Actualizar(ctx context.Context, gasto *Gasto) (*Gasto, error)
		Eliminar(ctx context.Context, id int) error
	}

	TipoPagoRepository interface {
		Listar(ctx context.Context) ([]*TipoPago, error)
	}

	GastoService interface {
		Listar(ctx context.Context, filtros GastoFiltros) ([]*Gasto, int, error)
		ObtenerGasto(ctx context.Context, id int, empresaID int) (*Gasto, error)
		ListarTiposPago(ctx context.Context) ([]*TipoPago, error)
		RegistrarGasto(ctx context.Context, gasto *Gasto) (*Gasto, error)
		ActualizarGasto(ctx context.Context, gasto *Gasto) (*Gasto, error)
		EliminarGasto(ctx context.Context, id int, empresaID int) error
	}
)
