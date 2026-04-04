package routes

import (
	"rentals-go/di"
	_ "rentals-go/docs"
	adminroutes "rentals-go/internal/routes/admin"
	userroutes "rentals-go/internal/routes/usuario"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Register(app *fiber.App, appDI *di.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "rentals-go listo",
		})
	})
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/catalogos/monedas", appDI.MonedaCtrl.Listar)
	app.Get("/catalogos/monedas/:codigo", appDI.MonedaCtrl.Obtener)

	adminroutes.Register(app, appDI)
	userroutes.Register(app, appDI)
}
