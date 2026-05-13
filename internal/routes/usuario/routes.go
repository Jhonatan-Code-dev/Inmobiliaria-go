package usuario

import (
	"rentals-go/di"
	"rentals-go/internal/routes/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appDI *di.App) {
	registrarRutasAuth(app.Group("/auth"), appDI)
	registrarRutasAuth(app.Group("/api/auth"), appDI)

	registrarRutaPerfil(app.Group("/me"), appDI)
	registrarRutaPerfil(app.Group("/api/me"), appDI)

	// Módulo de Gastos
	gastos := app.Group("/api/user/gastos")
	gastos.Use(middlewares.TenantAuth(appDI.Config))
	gastos.Get("/tipos-pago", appDI.GastoCtrl.ListarTiposPago)
	gastos.Get("/", appDI.GastoCtrl.Listar)
	gastos.Post("/", appDI.GastoCtrl.Crear)
	gastos.Put("/:id", appDI.GastoCtrl.Actualizar)
	gastos.Delete("/:id", appDI.GastoCtrl.Eliminar)
	gastos.Get("/reporte/excel", appDI.GastoCtrl.ExportarExcel)
	gastos.Get("/reporte/pdf", appDI.GastoCtrl.ExportarPDF)

	// Módulo de Clientes
	clientes := app.Group("/api/user/clientes")
	clientes.Use(middlewares.TenantAuth(appDI.Config))
	clientes.Get("/tipos-identificacion", appDI.ClienteCtrl.ListarTiposIdentificacion)
	clientes.Get("/", appDI.ClienteCtrl.Listar)
	clientes.Get("/:id", appDI.ClienteCtrl.Obtener)
	clientes.Post("/", appDI.ClienteCtrl.Crear)
	clientes.Put("/:id", appDI.ClienteCtrl.Actualizar)
	clientes.Delete("/:id", appDI.ClienteCtrl.Eliminar)

	// Módulo de Inmuebles
	inmuebles := app.Group("/api/user/inmuebles")
	inmuebles.Use(middlewares.TenantAuth(appDI.Config))
	inmuebles.Get("/", appDI.InmuebleCtrl.Listar)
	inmuebles.Get("/:id", appDI.InmuebleCtrl.Obtener)
	inmuebles.Post("/", appDI.InmuebleCtrl.Crear)
	inmuebles.Put("/:id", appDI.InmuebleCtrl.Actualizar)
	inmuebles.Delete("/:id", appDI.InmuebleCtrl.Eliminar)
	inmuebles.Get("/:id/unidades", appDI.InmuebleCtrl.ListarUnidades)
	inmuebles.Post("/:id/unidades", appDI.InmuebleCtrl.CrearUnidad)
	inmuebles.Get("/:id/unidades/:unidadId", appDI.InmuebleCtrl.ObtenerUnidad)
	inmuebles.Put("/:id/unidades/:unidadId", appDI.InmuebleCtrl.ActualizarUnidad)
	inmuebles.Delete("/:id/unidades/:unidadId", appDI.InmuebleCtrl.EliminarUnidad)

	// Módulo de Alquileres
	alquileres := app.Group("/api/user/alquileres")
	alquileres.Use(middlewares.TenantAuth(appDI.Config))
	alquileres.Get("/", appDI.AlquilerCtrl.Listar)
	alquileres.Get("/:id", appDI.AlquilerCtrl.Obtener)
	alquileres.Get("/activos/selector", appDI.AlquilerCtrl.ListarActivosSelector)
	alquileres.Post("/", appDI.AlquilerCtrl.Crear)
	alquileres.Put("/:id", appDI.AlquilerCtrl.Actualizar)
	alquileres.Delete("/:id", appDI.AlquilerCtrl.Eliminar)
	alquileres.Post("/:id/terminar", appDI.AlquilerCtrl.TerminarContrato)
	
	// Generación y Plantillas
	alquileres.Get("/plantillas", appDI.AlquilerCtrl.ListarPlantillas)
	alquileres.Post("/plantillas", appDI.AlquilerCtrl.GuardarPlantilla)
	alquileres.Delete("/plantillas/:id", appDI.AlquilerCtrl.EliminarPlantilla)
	alquileres.Get("/:id/generar-documento", appDI.AlquilerCtrl.GenerarDocumento)
	alquileres.Get("/:id/descargar-word", appDI.AlquilerCtrl.DescargarWord)
	alquileres.Post("/generar-borrador", appDI.AlquilerCtrl.GenerarBorrador)

	// Módulo de Pagos de alquiler
	pagos := app.Group("/api/user/pagos")
	pagos.Use(middlewares.TenantAuth(appDI.Config))
	pagos.Get("/", appDI.AlquilerCtrl.ListarPagos)
	pagos.Post("/", appDI.AlquilerCtrl.RegistrarPago)
	pagos.Get("/pendientes", appDI.AlquilerCtrl.PendientesPago)
	pagos.Get("/:id", appDI.AlquilerCtrl.ObtenerPago)
	pagos.Put("/:id", appDI.AlquilerCtrl.ActualizarPago)
	pagos.Delete("/:id", appDI.AlquilerCtrl.AnularPago)

	// Módulo de Staff
	staff := app.Group("/api/user/staff")
	staff.Use(middlewares.TenantAuth(appDI.Config))
	staff.Get("/roles", appDI.StaffCtrl.ListarRoles)
	staff.Get("/", appDI.StaffCtrl.Listar)
	staff.Get("/:id", appDI.StaffCtrl.Obtener)
	staff.Post("/", appDI.StaffCtrl.Crear)
	staff.Put("/:id", appDI.StaffCtrl.Actualizar)
	staff.Delete("/:id", appDI.StaffCtrl.Eliminar)

	// Módulo de Cargos
	cargos := app.Group("/api/user/cargos")
	cargos.Use(middlewares.TenantAuth(appDI.Config))
	cargos.Get("/", appDI.CargoCtrl.Listar)
	cargos.Get("/:id", appDI.CargoCtrl.Obtener)
	cargos.Post("/", appDI.CargoCtrl.Crear)
	cargos.Put("/:id", appDI.CargoCtrl.Actualizar)
	cargos.Delete("/:id", appDI.CargoCtrl.Eliminar)

	// Módulo de Servicios (Mediciones)
	servicios := app.Group("/api/user/servicios")
	servicios.Use(middlewares.TenantAuth(appDI.Config))
	servicios.Get("/", appDI.ServicioCtrl.Listar)
	servicios.Post("/", appDI.ServicioCtrl.Crear)
	servicios.Post("/masivo", appDI.ServicioCtrl.RegistrarMasivo)
	servicios.Post("/registrar-y-cobrar", appDI.ServicioCtrl.RegistrarYCobrar)
	servicios.Get("/ultimo/:contrato_id", appDI.ServicioCtrl.ObtenerUltima)
	servicios.Get("/pendientes", appDI.ServicioCtrl.ListarPendientes)
	servicios.Get("/:id", appDI.ServicioCtrl.Obtener)
	servicios.Put("/:id", appDI.ServicioCtrl.Actualizar)
	servicios.Delete("/:id", appDI.ServicioCtrl.Eliminar)

	// Módulo de Tickets (Mantenimiento)
	tickets := app.Group("/api/user/tickets")
	tickets.Use(middlewares.TenantAuth(appDI.Config))
	tickets.Get("/", appDI.TicketCtrl.Listar)
	tickets.Get("/resumen", appDI.TicketCtrl.Resumen)
	tickets.Get("/cola-trabajo", appDI.TicketCtrl.ColaTrabajo)
	tickets.Get("/config-formulario", appDI.TicketCtrl.ConfigFormulario)
	tickets.Post("/", appDI.TicketCtrl.Crear)
	tickets.Get("/:id", appDI.TicketCtrl.Obtener)
	tickets.Put("/:id", appDI.TicketCtrl.Actualizar)
	tickets.Patch("/:id/estado", appDI.TicketCtrl.CambiarEstado)
	tickets.Delete("/:id", appDI.TicketCtrl.Eliminar)

	// Módulo de Dashboard / KPIs
	dashboard := app.Group("/api/user/dashboard")
	dashboard.Use(middlewares.TenantAuth(appDI.Config))
	dashboard.Get("/", appDI.DashboardCtrl.ResumenGeneral)
	dashboard.Get("/ocupacion", appDI.DashboardCtrl.Ocupacion)
	dashboard.Get("/morosidad", appDI.DashboardCtrl.Morosidad)
	dashboard.Get("/financiero", appDI.DashboardCtrl.ReporteFinanciero)
	dashboard.Get("/contratos-por-vencer", appDI.DashboardCtrl.ContratosProximosVencer)
	dashboard.Get("/estado-cuenta/:clienteId", appDI.DashboardCtrl.EstadoCuentaCliente)
	dashboard.Get("/top-unidades", appDI.DashboardCtrl.TopUnidades)

	// Módulo de Asistencia
	asistencia := app.Group("/api/user/asistencia")
	asistencia.Use(middlewares.TenantAuth(appDI.Config))
	
	// Operaciones de Empleado
	asistencia.Post("/marcar", appDI.AsistenciaCtrl.MarcarAsistencia)
	asistencia.Get("/mi-historial", appDI.AsistenciaCtrl.MiHistorial)
	asistencia.Post("/permisos", appDI.AsistenciaCtrl.SolicitarPermiso)

	// Operaciones de Administrador
	asistencia.Get("/permisos", appDI.AsistenciaCtrl.ListarPermisos)
	asistencia.Get("/registros", appDI.AsistenciaCtrl.ListarRegistros)
	asistencia.Get("/reporte", appDI.AsistenciaCtrl.ConsultarReporteAsistencia)
	asistencia.Delete("/registros/:id", appDI.AsistenciaCtrl.EliminarRegistro)
	asistencia.Post("/horarios", appDI.AsistenciaCtrl.AsignarHorario)
	asistencia.Get("/horarios/detalle", appDI.AsistenciaCtrl.ObtenerHorario)
	asistencia.Put("/permisos/:id/estado", appDI.AsistenciaCtrl.DecidirPermiso)

	// Configuración Global
	asistencia.Get("/configuracion", appDI.AsistenciaCtrl.ObtenerConfiguracion)
	asistencia.Post("/configuracion", appDI.AsistenciaCtrl.ActualizarConfiguracion)

	// Exportaciones
	asistencia.Get("/reporte/excel", appDI.AsistenciaCtrl.ExportarReporteExcel)
	asistencia.Get("/reporte/pdf", appDI.AsistenciaCtrl.ExportarReportePDF)
}

func registrarRutasAuth(group fiber.Router, appDI *di.App) {
	group.Post("/login", appDI.UsuarioCtrl.Login)
	group.Post("/logout", appDI.UsuarioCtrl.Logout)
}

func registrarRutaPerfil(group fiber.Router, appDI *di.App) {
	group.Use(middlewares.TenantAuth(appDI.Config))
	group.Get("/", appDI.UsuarioCtrl.Perfil)
	group.Patch("/password", appDI.UsuarioCtrl.CambiarPassword)
}
