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
	alquileres.Post("/", appDI.AlquilerCtrl.Crear)

	// Módulo de Pagos de alquiler
	pagos := app.Group("/api/user/pagos")
	pagos.Use(middlewares.TenantAuth(appDI.Config))
	pagos.Post("/", appDI.AlquilerCtrl.RegistrarPago)
	pagos.Get("/pendientes", appDI.AlquilerCtrl.PendientesPago)
}

func registrarRutasAuth(group fiber.Router, appDI *di.App) {
	group.Post("/login", appDI.UsuarioCtrl.Login)
	group.Post("/logout", appDI.UsuarioCtrl.Logout)
}

func registrarRutaPerfil(group fiber.Router, appDI *di.App) {
	group.Use(middlewares.TenantAuth(appDI.Config))
	group.Get("/", appDI.UsuarioCtrl.Perfil)
}
