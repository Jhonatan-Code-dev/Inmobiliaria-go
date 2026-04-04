package main

import (
	"log"

	"rentals-go/di"
	"rentals-go/internal/routes"

	"github.com/gofiber/fiber/v2"
)

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
	routes.Register(app, appDI)
	log.Printf("API corriendo en http://localhost:%s\n", appDI.Config.Port)
	log.Fatal(app.Listen(":" + appDI.Config.Port))
}
