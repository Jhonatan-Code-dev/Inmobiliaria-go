package controller

import (
	"net/http"
	"time"

	"rentals-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type usuarioLoginRequest struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type usuarioLoginResponse struct {
	Token     string           `json:"token"`
	EmpresaID int              `json:"empresa_id"`
	User      usuarioResponse  `json:"user"`
	Emp       *empresaResponse `json:"empresa,omitempty"`
}

type UsuarioController struct {
	svc *service.UsuarioService
}

func NewUsuarioController(svc *service.UsuarioService) *UsuarioController {
	return &UsuarioController{svc: svc}
}

// Login godoc
// @Summary Iniciar sesión como usuario
// @Description Permite a un usuario autenticarse y obtener un token JWT. Además, establece una cookie de sesión.
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param request body usuarioLoginRequest true "Credenciales de usuario"
// @Success 200 {object} usuarioLoginResponse
// @Failure 401 {object} errorResponse "Credenciales inválidas"
// @Router /auth/login [post]
func (h *UsuarioController) Login(c *fiber.Ctx) error {
	var req usuarioLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	token, user, emp, err := h.svc.Login(c.Context(), req.Usuario, req.Contrasena)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(errorResponse{Message: "credenciales inválidas"})
	}
	c.Cookie(&fiber.Cookie{Name: "token_usuario", Value: token, HTTPOnly: true, Path: "/"})
	return c.JSON(usuarioLoginResponse{
		Token:     token,
		EmpresaID: user.EmpresaID,
		User: usuarioResponse{
			ID:        user.ID,
			Usuario:   user.Usuario,
			EmpresaID: user.EmpresaID,
		},
		Emp: mapEmpresaResponse(emp),
	})
}

// Logout godoc
// @Summary Cerrar sesión
// @Description Elimina el token de sesión (cookie) del navegador.
// @Tags Usuarios
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *UsuarioController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token_usuario",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})
	return c.JSON(fiber.Map{"message": "sesión cerrada"})
}

// Perfil godoc
// @Summary Obtener perfil del usuario actual
// @Description Retorna los datos del usuario autenticado y su empresa.
// @Tags Usuarios
// @Security ApiKeyAuth
// @Success 200 {object} usuarioLoginResponse
// @Failure 401 {object} errorResponse "No autorizado"
// @Router /me [get]
func (h *UsuarioController) Perfil(c *fiber.Ctx) error {
	idVal := c.Locals("usuario_id")
	if idVal == nil {
		return fiber.ErrUnauthorized
	}
	id := idVal.(int)
	user, emp, err := h.svc.Perfil(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: "error interno"})
	}
	return c.JSON(usuarioLoginResponse{
		EmpresaID: user.EmpresaID,
		User: usuarioResponse{
			ID:        user.ID,
			Usuario:   user.Usuario,
			EmpresaID: user.EmpresaID,
		},
		Emp: mapEmpresaResponse(emp),
	})
}
