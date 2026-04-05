package admin

import (
	"rentals-go/di"
	"rentals-go/internal/routes/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appDI *di.App) {
	group := app.Group("/admin")
	group.Post("/login", appDI.AdminCtrl.Login)

	protected := group.Use(middlewares.AdminAuth(appDI.Config))
	protected.Get("/me", appDI.AdminCtrl.Perfil)
	protected.Patch("/credenciales", appDI.AdminCtrl.ActualizarCredenciales)
	protected.Get("/empresas", appDI.AdminCtrl.ListarEmpresas)
	protected.Post("/empresas", appDI.AdminCtrl.CrearEmpresa)
	protected.Get("/empresas/:id", appDI.AdminCtrl.ObtenerEmpresa)
	protected.Put("/empresas/:id", appDI.AdminCtrl.ActualizarEmpresa)
	protected.Delete("/empresas/:id", appDI.AdminCtrl.EliminarEmpresa)
}
