package usuario

import (
	"rentals-go/di"
	"rentals-go/internal/routes/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appDI *di.App) {
	authGroup := app.Group("/auth")
	authGroup.Post("/login", appDI.UsuarioCtrl.Login)

	me := app.Group("/me")
	me.Use(middlewares.TenantAuth(appDI.Config))
	me.Get("/", appDI.UsuarioCtrl.Perfil)
}
