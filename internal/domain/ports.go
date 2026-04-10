package domain

import (
	"context"
	"time"
)

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

	TipoIdentificacionRepository interface {
		ListarActivos(ctx context.Context) ([]*TipoIdentificacion, error)
		ExisteActivo(ctx context.Context, id int) (bool, error)
	}

	ClienteRepository interface {
		ListarPaginado(ctx context.Context, filtros ClienteFiltros) ([]*Cliente, int, error)
		BuscarPorID(ctx context.Context, id int) (*Cliente, error)
		Crear(ctx context.Context, c *Cliente) (*Cliente, error)
		Actualizar(ctx context.Context, c *Cliente) (*Cliente, error)
		Eliminar(ctx context.Context, id int) error
	}

	InmuebleRepository interface {
		ListarPaginado(ctx context.Context, filtros InmuebleFiltros) ([]*Inmueble, int, error)
		BuscarPorID(ctx context.Context, id int) (*Inmueble, error)
		Crear(ctx context.Context, inmueble *Inmueble) (*Inmueble, error)
		Actualizar(ctx context.Context, inmueble *Inmueble) (*Inmueble, error)
		Eliminar(ctx context.Context, id int) error
		ListarUnidades(ctx context.Context, propiedadID int) ([]*Unidad, error)
		BuscarUnidadPorID(ctx context.Context, id int) (*Unidad, error)
		CrearUnidad(ctx context.Context, unidad *Unidad) (*Unidad, error)
		ActualizarUnidad(ctx context.Context, unidad *Unidad) (*Unidad, error)
		EliminarUnidad(ctx context.Context, id int) error
	}

	AlquilerRepository interface {
		ListarPaginado(ctx context.Context, filtros AlquilerFiltros) ([]*Alquiler, int, error)
		BuscarPorID(ctx context.Context, id int) (*Alquiler, error)
		Crear(ctx context.Context, alquiler *Alquiler) (*Alquiler, error)
	}

	PagoAlquilerRepository interface {
		Registrar(ctx context.Context, pago *RegistroPagoAlquiler) (*PagoAlquiler, error)
		ListarPendientesMesActual(ctx context.Context, empresaID int, now time.Time) ([]*PagoPendiente, error)
	}

	GastoService interface {
		Listar(ctx context.Context, filtros GastoFiltros) ([]*Gasto, int, error)
		ObtenerGasto(ctx context.Context, id int, empresaID int) (*Gasto, error)
		ListarTiposPago(ctx context.Context) ([]*TipoPago, error)
		RegistrarGasto(ctx context.Context, gasto *Gasto) (*Gasto, error)
		ActualizarGasto(ctx context.Context, gasto *Gasto) (*Gasto, error)
		EliminarGasto(ctx context.Context, id int, empresaID int) error
	}

	ClienteService interface {
		Listar(ctx context.Context, filtros ClienteFiltros) ([]*Cliente, int, error)
		ObtenerCliente(ctx context.Context, id int, empresaID int) (*Cliente, error)
		ListarTiposIdentificacion(ctx context.Context) ([]*TipoIdentificacion, error)
		RegistrarCliente(ctx context.Context, c *Cliente) (*Cliente, error)
		ActualizarCliente(ctx context.Context, c *Cliente) (*Cliente, error)
		EliminarCliente(ctx context.Context, id int, empresaID int) error
	}

	InmuebleService interface {
		Listar(ctx context.Context, filtros InmuebleFiltros) ([]*Inmueble, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*Inmueble, error)
		Crear(ctx context.Context, inmueble *Inmueble) (*Inmueble, error)
		Actualizar(ctx context.Context, inmueble *Inmueble) (*Inmueble, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
		ListarUnidades(ctx context.Context, propiedadID int, empresaID int) ([]*Unidad, error)
		ObtenerUnidad(ctx context.Context, propiedadID int, unidadID int, empresaID int) (*Unidad, error)
		CrearUnidad(ctx context.Context, propiedadID int, empresaID int, unidad *Unidad) (*Unidad, error)
		ActualizarUnidad(ctx context.Context, propiedadID int, empresaID int, unidad *Unidad) (*Unidad, error)
		EliminarUnidad(ctx context.Context, propiedadID int, unidadID int, empresaID int) error
	}

	AlquilerService interface {
		Listar(ctx context.Context, filtros AlquilerFiltros) ([]*Alquiler, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*Alquiler, error)
		Crear(ctx context.Context, alquiler *Alquiler) (*Alquiler, error)
	}

	PagoAlquilerService interface {
		Registrar(ctx context.Context, pago *RegistroPagoAlquiler) (*PagoAlquiler, error)
		ListarPendientesMesActual(ctx context.Context, empresaID int) ([]*PagoPendiente, error)
	}
)
