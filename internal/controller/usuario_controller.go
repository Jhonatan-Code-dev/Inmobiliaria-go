package controller

import (
	"net/http"

	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/moneda"
	"rentals-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

type usuarioLoginRequest struct {
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena"`
}

type usuarioLoginResponse struct {
	Token string           `json:"token"`
	User  usuarioResponse  `json:"user"`
	Emp   *empresaResponse `json:"empresa,omitempty"`
}

type usuarioResponse struct {
	ID        int    `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos,omitempty"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono,omitempty"`
}

type UsuarioController struct {
	svc *service.UsuarioService
}

func NewUsuarioController(svc *service.UsuarioService) *UsuarioController {
	return &UsuarioController{svc: svc}
}

// LoginUsuario godoc
// @Summary Login de usuario
// @Tags usuario
// @Accept json
// @Produce json
// @Param credentials body usuarioLoginRequest true "Credenciales"
// @Success 200 {object} usuarioLoginResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/login [post]
func (h *UsuarioController) Login(c *fiber.Ctx) error {
	var req usuarioLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	token, user, emp, err := h.svc.Login(c.Context(), req.Correo, req.Contrasena)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "credenciales inválidas")
	}
	c.Cookie(&fiber.Cookie{Name: "token_usuario", Value: token, HTTPOnly: true, Path: "/"})
	return c.JSON(usuarioLoginResponse{
		Token: token,
		User: usuarioResponse{
			ID:        user.ID,
			Nombres:   user.Nombres,
			Apellidos: user.Apellidos,
			Correo:    user.Correo,
			Telefono:  user.Telefono,
		},
		Emp: mapEmpresa(emp),
	})
}

// Perfil godoc
// @Summary Perfil del usuario autenticado
// @Tags usuario
// @Produce json
// @Success 200 {object} usuarioLoginResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /me [get]
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
			ID:        user.ID,
			Nombres:   user.Nombres,
			Apellidos: user.Apellidos,
			Correo:    user.Correo,
			Telefono:  user.Telefono,
		},
		Emp: mapEmpresa(emp),
	})
}

func mapEmpresa(emp *domain.Empresa) *empresaResponse {
	if emp == nil {
		return nil
	}
	resp := &empresaResponse{
		ID:              emp.ID,
		Nombre:          emp.Nombre,
		DocumentoFiscal: emp.DocumentoFiscal,
		Correo:          emp.Correo,
		Telefono:        emp.Telefono,
		Direccion:       emp.Direccion,
		Ciudad:          emp.Ciudad,
		Pais:            emp.Pais,
		Moneda:          emp.Moneda,
		MaximoUsuarios:  emp.MaximoUsuarios,
		Estado:          emp.Estado,
	}
	if info, err := moneda.ObtenerInfo(emp.Moneda); err == nil {
		mapped := monedaResponse{
			Codigo:     info.Codigo,
			Decimales:  info.Decimales,
			Incremento: info.Incremento,
			Regiones:   nil,
			Render: monedaRenderResponse{
				Metodo:                info.Render.Metodo,
				Currency:              info.Render.Currency,
				MinimumFractionDigits: info.Render.MinimumFractionDigits,
				MaximumFractionDigits: info.Render.MaximumFractionDigits,
			},
		}
		resp.MonedaInfo = &mapped
	}
	return resp
}
