package middlewares

import (
	"strings"

	"rentals-go/config/env"
	"rentals-go/internal/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

const (
	rolAdmin  = "admin"
	rolTenant = "usuario"
)

// AdminAuth valida token de administrador.
func AdminAuth(cfg *env.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := parseToken(c, []byte(cfg.JWTSecret))
		if err != nil || claims.AdminID == 0 || claims.Rol != rolAdmin {
			return fiber.ErrUnauthorized
		}
		c.Locals("admin_id", claims.AdminID)
		return c.Next()
	}
}

// TenantAuth valida token de usuario.
func TenantAuth(cfg *env.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := parseToken(c, []byte(cfg.JWTSecret))
		if err != nil || claims.UsuarioID == 0 || claims.Rol != rolTenant {
			return fiber.ErrUnauthorized
		}
		c.Locals("usuario_id", claims.UsuarioID)
		c.Locals("empresa_id", claims.EmpresaID)
		return c.Next()
	}
}

func parseToken(c *fiber.Ctx, secret []byte) (*auth.Claims, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fiber.ErrUnauthorized
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, fiber.ErrUnauthorized
	}
	return auth.ValidarToken(parts[1], secret)
}
