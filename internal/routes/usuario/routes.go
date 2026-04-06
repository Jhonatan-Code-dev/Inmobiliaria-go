package usuario

import (
	"rentals-go/di"
	"rentals-go/internal/routes/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appDI *di.App) {
	// Rutas de autenticación pública
	authGroup := app.Group("/auth")
	authGroup.Post("/login", appDI.UsuarioCtrl.Login)
	authGroup.Post("/logout", appDI.UsuarioCtrl.Logout)

	// Perfil del usuario autenticado
	me := app.Group("/me")
	me.Use(middlewares.TenantAuth(appDI.Config))
	me.Get("/", appDI.UsuarioCtrl.Perfil)
}
