package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"rentals-go/config/security"
	"rentals-go/internal/domain"
	"rentals-go/internal/service"

	"github.com/gofiber/fiber/v2"
)

// DTOs locales (aislados al controller).
type adminLoginRequest struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type adminLoginResponse struct {
	Token string `json:"token"`
}

type crearEmpresaRequest struct {
	Empresa struct {
		Nombre string `json:"nombre"`
		Pais   string `json:"pais"`
		Moneda string `json:"moneda"`
	} `json:"empresa"`
	Usuario struct {
		Username string `json:"usuario"`
		Password string `json:"password"`
	} `json:"usuario"`
}

type crearEmpresaResponse struct {
	EmpresaID int `json:"empresa_id"`
	UsuarioID int `json:"usuario_id"`
}

type actualizarEmpresaRequest struct {
	Nombre string `json:"nombre"`
	Pais   string `json:"pais"`
	Moneda string `json:"moneda"`
	Estado string `json:"estado"`
}

type adminCredencialesRequest struct {
	ID         int    `json:"id"`
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type adminCredencialesResponse struct {
	ID      int    `json:"id"`
	Usuario string `json:"usuario"`
}

type adminProfileResponse struct {
	ID      int    `json:"id"`
	Usuario string `json:"usuario"`
	Activo  bool   `json:"activo"`
}

type AdminController struct {
	svc *service.AdminService
}

func NewAdminController(svc *service.AdminService) *AdminController {
	return &AdminController{svc: svc}
}

// LoginAdmin godoc
// @Summary Login admin
// @Description Autentica a un administrador y devuelve el JWT que el frontend debe enviar en el header Authorization.
// @Tags admin
// @Accept json
// @Produce json
// @Param credentials body adminLoginRequest true "Credenciales"
// @Success 200 {object} adminLoginResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/login [post]
func (h *AdminController) Login(c *fiber.Ctx) error {
	var req adminLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	token, err := h.svc.Login(c.Context(), req.Usuario, req.Contrasena)
	if err != nil {
		return fiber.NewError(http.StatusUnauthorized, "credenciales inválidas")
	}
	c.Cookie(&fiber.Cookie{Name: "token_admin", Value: token, HTTPOnly: true, Path: "/"})
	return c.JSON(adminLoginResponse{Token: token})
}

// LogoutAdmin godoc
// @Summary Logout admin
// @Description Cierra la sesión del administrador actual eliminando la cookie de autenticación.
// @Tags admin
// @Produce json
// @Success 200 {object} map[string]string
// @Router /admin/logout [post]
func (h *AdminController) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token_admin",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})
	return c.JSON(fiber.Map{"message": "sesión cerrada"})
}

// PerfilAdmin godoc
// @Summary Perfil del administrador autenticado
// @Description Retorna los datos del administrador autenticado a partir del token Bearer enviado por el frontend.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} adminProfileResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/me [get]
func (h *AdminController) Perfil(c *fiber.Ctx) error {
	adminIDVal := c.Locals("admin_id")
	if adminIDVal == nil {
		return fiber.ErrUnauthorized
	}
	adminID, ok := adminIDVal.(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	adminActual, err := h.svc.Perfil(c.Context(), adminID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(adminProfileResponse{
		ID:      adminActual.ID,
		Usuario: adminActual.Usuario,
		Activo:  adminActual.Activo,
	})
}

// CrearEmpresa godoc
// @Summary Alta de empresa + usuario principal
// @Description Crea una empresa y su usuario principal en una sola operacion. Requiere autenticacion de administrador. Solo nombre, pais, usuario y contraseña.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body crearEmpresaRequest true "Datos de empresa y usuario"
// @Success 201 {object} crearEmpresaResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/empresas [post]
func (h *AdminController) CrearEmpresa(c *fiber.Ctx) error {
	var req crearEmpresaRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validarCrearEmpresaRequest(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	hasher := security.NewServicioHash()
	hash, err := hasher.Encriptar(req.Usuario.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	emp := &domain.Empresa{
		Nombre: req.Empresa.Nombre,
		Pais:   req.Empresa.Pais,
		Moneda: defaultString(req.Empresa.Moneda, "PEN"),
	}
	u := &domain.Usuario{
		Usuario:        req.Usuario.Username,
		HashContrasena: hash,
	}
	createdEmp, createdUser, err := h.svc.CrearEmpresaConUsuario(c.Context(), emp, u, 0)
	if err != nil {
		if err == service.ErrMonedaInvalida {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusCreated).JSON(crearEmpresaResponse{
		EmpresaID: createdEmp.ID,
		UsuarioID: createdUser.ID,
	})
}

// ObtenerEmpresa godoc
// @Summary Obtener empresa por ID
// @Description Devuelve el detalle de una empresa.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de empresa"
// @Success 200 {object} empresaResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/empresas/{id} [get]
func (h *AdminController) ObtenerEmpresa(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return fiber.ErrBadRequest
	}
	empresa, err := h.svc.ObtenerEmpresa(c.Context(), id)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(mapEmpresaResponse(empresa))
}

// ListarEmpresas godoc
// @Summary Listado de empresas
// @Description Retorna todas las empresas registradas.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} empresaResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/empresas [get]
func (h *AdminController) ListarEmpresas(c *fiber.Ctx) error {
	list, err := h.svc.ListarEmpresas(c.Context())
	if err != nil {
		return fiber.ErrInternalServerError
	}
	resp := make([]empresaResponse, 0, len(list))
	for _, e := range list {
		mapped := mapEmpresaResponse(e)
		if mapped != nil {
			resp = append(resp, *mapped)
		}
	}
	return c.JSON(resp)
}

// ActualizarEmpresa godoc
// @Summary Actualizar empresa
// @Description Actualiza los datos generales de una empresa.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de empresa"
// @Param body body actualizarEmpresaRequest true "Datos de empresa"
// @Success 200 {object} empresaResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/empresas/{id} [put]
func (h *AdminController) ActualizarEmpresa(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var req actualizarEmpresaRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if req.Nombre == "" {
		return fiber.ErrBadRequest
	}

	empresa, err := h.svc.ActualizarEmpresa(c.Context(), &domain.Empresa{
		ID:     id,
		Nombre: req.Nombre,
		Pais:   req.Pais,
		Moneda: defaultString(req.Moneda, "PEN"),
		Estado: req.Estado,
	})
	if err != nil {
		if err == service.ErrMonedaInvalida {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(mapEmpresaResponse(empresa))
}

// EliminarEmpresa godoc
// @Summary Eliminar empresa
// @Description Elimina una empresa por ID.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de empresa"
// @Success 204
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/empresas/{id} [delete]
func (h *AdminController) EliminarEmpresa(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := h.svc.EliminarEmpresa(c.Context(), id); err != nil {
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(http.StatusNoContent)
}

// ActualizarCredenciales godoc
// @Summary Actualizar credenciales
// @Description Actualiza el usuario y contraseña de un administrador. Si no se envía ID en el body, se actualiza el admin autenticado.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body adminCredencialesRequest true "Datos de credenciales"
// @Success 200 {object} adminCredencialesResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /admin/credenciales [patch]
func (h *AdminController) ActualizarCredenciales(c *fiber.Ctx) error {
	adminIDVal := c.Locals("admin_id")
	if adminIDVal == nil {
		return fiber.ErrUnauthorized
	}
	authAdminID := adminIDVal.(int)

	var req adminCredencialesRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// Si el request trae un ID, usamos ese, si no, el del token.
	targetID := req.ID
	if targetID == 0 {
		targetID = authAdminID
	}

	adminActualizado, err := h.svc.ActualizarCredenciales(c.Context(), targetID, req.Usuario, req.Contrasena)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
			return fiber.NewError(http.StatusBadRequest, "el usuario ya está en uso")
		}
		if err == service.ErrUsuarioAdminInvalido || err == service.ErrContrasenaInvalida || strings.Contains(err.Error(), "obligatorio") {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(adminCredencialesResponse{
		ID:      adminActualizado.ID,
		Usuario: adminActualizado.Usuario,
	})
}

func defaultString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func validarCrearEmpresaRequest(req crearEmpresaRequest) error {
	if strings.TrimSpace(req.Empresa.Nombre) == "" {
		return fiber.NewError(http.StatusBadRequest, "empresa.nombre es obligatorio")
	}
	if req.Empresa.Pais != "" && len(strings.TrimSpace(req.Empresa.Pais)) != 2 {
		return fiber.NewError(http.StatusBadRequest, "empresa.pais debe ser un codigo ISO de 2 letras, por ejemplo PE")
	}
	if strings.TrimSpace(req.Usuario.Username) == "" {
		return fiber.NewError(http.StatusBadRequest, "usuario es obligatorio")
	}
	if strings.TrimSpace(req.Usuario.Password) == "" {
		return fiber.NewError(http.StatusBadRequest, "usuario.password es obligatorio")
	}
	return nil
}

func parseIDParam(c *fiber.Ctx) (int, error) {
	return strconv.Atoi(c.Params("id"))
}
