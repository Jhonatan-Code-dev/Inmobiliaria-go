// Package security proporciona utilidades para manejo de JWT, cookies y hashing
package security

import (
	"time"

	"rentals-go/internal/pkg/tiempo"

	"github.com/gofiber/fiber/v2"
)

// CookieConfig define la configuración básica de la cookie desde el entorno
type CookieConfig struct {
	MaxAge time.Duration
	Secure bool
}

// Nombres de cookies constantes para evitar colisiones entre sistemas
const (
	CookieNameAdmin = "admin_access_token"
	CookieNameUser  = "user_access_token"
)

// GestorCookie centraliza la lógica de creación y eliminación de cookies
type GestorCookie struct {
	config CookieConfig
}

// NewGestorCookie crea una nueva instancia del gestor
func NewGestorCookie(cfg CookieConfig) *GestorCookie {
	return &GestorCookie{config: cfg}
}

// Crear coloca una cookie de forma segura (HttpOnly y SameSite por defecto)
func (g *GestorCookie) Crear(ctx *fiber.Ctx, nombre, valor string) {
	ctx.Cookie(&fiber.Cookie{
		Name:     nombre,
		Value:    valor,
		Expires:  tiempo.AhoraUTC().Add(g.config.MaxAge),
		HTTPOnly: true, // Protege contra XSS
		Secure:   g.config.Secure,
		SameSite: "Lax", // Protege contra CSRF
		Path:     "/",
	})
}

// Eliminar invalida la cookie expirándola de inmediato
func (g *GestorCookie) Eliminar(ctx *fiber.Ctx, nombre string) {
	ctx.Cookie(&fiber.Cookie{
		Name:  nombre,
		Value: "",
		// Expiración en el pasado para que el navegador la borre
		Expires:  tiempo.AhoraUTC().Add(-24 * time.Hour),
		HTTPOnly: true,
		Secure:   g.config.Secure,
		Path:     "/",
	})
}

// SetTokenCookieAdmin coloca la cookie para el sistema administrativo
func (g *GestorCookie) SetTokenCookieAdmin(ctx *fiber.Ctx, token string) {
	g.Crear(ctx, CookieNameAdmin, token)
}

// SetTokenCookieUsuario coloca la cookie para el sistema de usuarios (empresa)
func (g *GestorCookie) SetTokenCookieUsuario(ctx *fiber.Ctx, token string) {
	g.Crear(ctx, CookieNameUser, token)
}
