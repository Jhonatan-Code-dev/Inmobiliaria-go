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
	origins := appDI.Config.AllowedOrigins
	if origins == "*" {
		// Fiber no permite "*" con AllowCredentials=true. 
		// Si es desarrollo y se quiere permitir todo, es común usar un string vacío o manejarlo por request, 
		// pero una forma común en este repo es listar los orígenes o dejar que el middleware maneje el "*" si credentials es false.
		// Para soportar cookies de cualquier origen en dev, usamos un pequeño truco o listamos los más comunes.
		origins = "http://localhost:3000,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	routes.Register(app, appDI)
	log.Printf("API corriendo en http://localhost:%s\n", appDI.Config.Port)
	log.Fatal(app.Listen(":" + appDI.Config.Port))
}
