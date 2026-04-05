package main

import (
	"log"

	"rentals-go/di"
	"rentals-go/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Rentals Go API
// @version 1.0
// @description API para autenticacion, administracion de empresas y catalogos de soporte para el frontend de Rentals Go.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT en formato: Bearer <token>
func main() {
	appDI, err := di.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if appDI.EntClient != nil {
			_ = appDI.EntClient.Close()
		}
	}()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: appDI.Config.AllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	routes.Register(app, appDI)
	log.Printf("API corriendo en http://localhost:%s\n", appDI.Config.Port)
	log.Fatal(app.Listen(":" + appDI.Config.Port))
}
