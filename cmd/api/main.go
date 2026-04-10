package main

import (
	"log"
	"net/http"
	"strings"

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

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := http.StatusInternalServerError
			message := "internal server error"

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			} else if err != nil {
				message = err.Error()
			}

			return c.Status(code).JSON(fiber.Map{
				"message": message,
			})
		},
	})
	origins := normalizarAllowedOrigins(appDI.Config.AllowedOrigins)

	corsConfig := cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}

	if origins == "*" {
		corsConfig.AllowOriginsFunc = func(origin string) bool {
			return strings.TrimSpace(origin) != ""
		}
	} else {
		corsConfig.AllowOrigins = origins
	}

	app.Use(cors.New(corsConfig))
	routes.Register(app, appDI)
	log.Printf("API corriendo en http://localhost:%s\n", appDI.Config.Port)
	log.Fatal(app.Listen(":" + appDI.Config.Port))
}

func normalizarAllowedOrigins(origins string) string {
	origins = strings.TrimSpace(origins)
	if origins == "" || origins == "*" {
		return "*"
	}

	partes := strings.Split(origins, ",")
	normalizados := make([]string, 0, len(partes))
	for _, origen := range partes {
		origen = strings.TrimSpace(origen)
		if origen == "" {
			continue
		}
		normalizados = append(normalizados, origen)
	}

	if len(normalizados) == 0 {
		return "*"
	}

	return strings.Join(normalizados, ",")
}
