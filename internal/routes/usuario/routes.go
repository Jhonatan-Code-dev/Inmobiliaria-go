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
}

func registrarRutasAuth(group fiber.Router, appDI *di.App) {
	group.Post("/login", appDI.UsuarioCtrl.Login)
	group.Post("/logout", appDI.UsuarioCtrl.Logout)
}

func registrarRutaPerfil(group fiber.Router, appDI *di.App) {
	group.Use(middlewares.TenantAuth(appDI.Config))
	group.Get("/", appDI.UsuarioCtrl.Perfil)
}
