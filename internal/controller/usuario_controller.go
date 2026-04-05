package controller

import (
	"net/http"

	"rentals-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type usuarioLoginRequest struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type usuarioLoginResponse struct {
	Token string           `json:"token"`
	User  usuarioResponse  `json:"user"`
	Emp   *empresaResponse `json:"empresa,omitempty"`
}

type UsuarioController struct {
	svc *service.UsuarioService
}

func NewUsuarioController(svc *service.UsuarioService) *UsuarioController {
	return &UsuarioController{svc: svc}
}

// Login godoc
func (h *UsuarioController) Login(c *fiber.Ctx) error {
	var req usuarioLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	token, user, emp, err := h.svc.Login(c.Context(), req.Usuario, req.Contrasena)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "credenciales inválidas")
	}
	c.Cookie(&fiber.Cookie{Name: "token_usuario", Value: token, HTTPOnly: true, Path: "/"})
	return c.JSON(usuarioLoginResponse{
		Token: token,
		User: usuarioResponse{
			ID:      user.ID,
			Usuario: user.Usuario,
		},
		Emp: mapEmpresaResponse(emp),
	})
}

// Perfil godoc
func (h *UsuarioController) Perfil(c *fiber.Ctx) error {
	idVal := c.Locals("usuario_id")
	if idVal == nil {
		return fiber.ErrUnauthorized
	}
	id := idVal.(int)
	user, emp, err := h.svc.Perfil(c.Context(), id)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(usuarioLoginResponse{
		User: usuarioResponse{
			ID:      user.ID,
			Usuario: user.Usuario,
		},
		Emp: mapEmpresaResponse(emp),
	})
}
