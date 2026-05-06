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
		Listar(ctx context.Context) ([]*Rol, error)
	}

	UsuarioRepository interface {
		Crear(ctx context.Context, u *Usuario) (*Usuario, error)
		BuscarPorUsuario(ctx context.Context, usuario string) (*Usuario, error)
		BuscarPerfil(ctx context.Context, id int) (*Usuario, *Empresa, error)
		ActualizarPassword(ctx context.Context, id int, hashContrasena string) error
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
		Actualizar(ctx context.Context, alquiler *Alquiler) (*Alquiler, error)
		Eliminar(ctx context.Context, id int) error
	}

	PagoAlquilerRepository interface {
		Registrar(ctx context.Context, pago *RegistroPagoAlquiler) (*PagoAlquiler, error)
		ListarPendientesMesActual(ctx context.Context, empresaID int, now time.Time) ([]*PagoPendiente, error)
		Listar(ctx context.Context, filtros PagoFiltros) ([]*PagoAlquiler, int, error)
		BuscarPorID(ctx context.Context, id int, empresaID int) (*PagoAlquiler, error)
		Actualizar(ctx context.Context, pago *PagoAlquiler) (*PagoAlquiler, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
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
		Actualizar(ctx context.Context, id int, empresaID int, alq *Alquiler) (*Alquiler, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
		Terminar(ctx context.Context, id int, empresaID int) error
	}

	PagoAlquilerService interface {
		Registrar(ctx context.Context, pago *RegistroPagoAlquiler) (*PagoAlquiler, error)
		ListarPendientesMesActual(ctx context.Context, empresaID int) ([]*PagoPendiente, error)
		ListarHistorial(ctx context.Context, filtros PagoFiltros) ([]*PagoAlquiler, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*PagoAlquiler, error)
		Actualizar(ctx context.Context, id int, empresaID int, notas *string, metodoPago string) (*PagoAlquiler, error)
		Anular(ctx context.Context, id int, empresaID int) error
	}

	StaffRepository interface {
		Listar(ctx context.Context, filtros StaffFiltros) ([]*Staff, int, error)
		BuscarPorID(ctx context.Context, id int, empresaID int) (*Staff, error)
		Crear(ctx context.Context, s *RegistroStaff, hash string) (*Staff, error)
		Actualizar(ctx context.Context, id int, empresaID int, rolID int, estado string) (*Staff, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}

	StaffService interface {
		Listar(ctx context.Context, filtros StaffFiltros) ([]*Staff, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*Staff, error)
		Registrar(ctx context.Context, s *RegistroStaff) (*Staff, error)
		Actualizar(ctx context.Context, id int, empresaID int, rolID int, estado string) (*Staff, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
		ListarRoles(ctx context.Context) ([]*Rol, error)
	}

	CargoRepository interface {
		Listar(ctx context.Context, filtros CargoFiltros) ([]*Cargo, int, error)
		BuscarPorID(ctx context.Context, id int, empresaID int) (*Cargo, error)
		Crear(ctx context.Context, c *Cargo) (*Cargo, error)
		Actualizar(ctx context.Context, c *Cargo) (*Cargo, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}

	CargoService interface {
		Listar(ctx context.Context, filtros CargoFiltros) ([]*Cargo, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*Cargo, error)
		Crear(ctx context.Context, c *RegistroCargo, empresaID int) (*Cargo, error)
		Actualizar(ctx context.Context, id int, empresaID int, c *RegistroCargo) (*Cargo, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}

	ServicioMedicionRepository interface {
		Listar(ctx context.Context, filtros ServicioMedicionFiltros) ([]*ServicioMedicion, int, error)
		BuscarPorID(ctx context.Context, id int, empresaID int) (*ServicioMedicion, error)
		Crear(ctx context.Context, s *ServicioMedicion) (*ServicioMedicion, error)
		Actualizar(ctx context.Context, s *ServicioMedicion) (*ServicioMedicion, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
		ObtenerUltimaLectura(ctx context.Context, contratoID int, tipo string) (*ServicioMedicion, error)
	}

	ServicioMedicionService interface {
		Listar(ctx context.Context, filtros ServicioMedicionFiltros) ([]*ServicioMedicion, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*ServicioMedicion, error)
		Registrar(ctx context.Context, r *RegistroLectura, empresaID int) (*ServicioMedicion, error)
		Actualizar(ctx context.Context, id int, empresaID int, lecturaActual float64) (*ServicioMedicion, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}

	TicketRepository interface {
		Listar(ctx context.Context, filtros TicketFiltros) ([]*Ticket, int, error)
		BuscarPorID(ctx context.Context, id int, empresaID int) (*Ticket, error)
		Crear(ctx context.Context, t *Ticket) (*Ticket, error)
		Actualizar(ctx context.Context, t *Ticket) (*Ticket, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}

	TicketService interface {
		Listar(ctx context.Context, filtros TicketFiltros) ([]*Ticket, int, error)
		Obtener(ctx context.Context, id int, empresaID int) (*Ticket, error)
		Crear(ctx context.Context, r *RegistroTicket, empresaID int) (*Ticket, error)
		Actualizar(ctx context.Context, id int, empresaID int, r *RegistroTicket, estado string) (*Ticket, error)
		Eliminar(ctx context.Context, id int, empresaID int) error
	}
)
