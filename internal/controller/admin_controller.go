package controller

import (
	"net/http"
	"strconv"
	"strings"

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

type errorResponse struct {
	Message string `json:"message"`
}

type crearEmpresaRequest struct {
	Empresa struct {
		Nombre          string `json:"nombre"`
		DocumentoFiscal string `json:"documento_fiscal"`
		Correo          string `json:"correo"`
		Telefono        string `json:"telefono"`
		Direccion       string `json:"direccion"`
		Ciudad          string `json:"ciudad"`
		Pais            string `json:"pais"`
		Moneda          string `json:"moneda"`
		MaximoUsuarios  int    `json:"maximo_usuarios"`
		Estado          string `json:"estado"`
	} `json:"empresa"`
	Usuario struct {
		Nombres   string `json:"nombres"`
		Apellidos string `json:"apellidos"`
		Correo    string `json:"correo"`
		Telefono  string `json:"telefono"`
		Password  string `json:"password"`
	} `json:"usuario"`
}

type crearEmpresaResponse struct {
	EmpresaID int `json:"empresa_id"`
	UsuarioID int `json:"usuario_id"`
}

type actualizarEmpresaRequest struct {
	Nombre          string `json:"nombre"`
	DocumentoFiscal string `json:"documento_fiscal"`
	Correo          string `json:"correo"`
	Telefono        string `json:"telefono"`
	Direccion       string `json:"direccion"`
	Ciudad          string `json:"ciudad"`
	Pais            string `json:"pais"`
	Moneda          string `json:"moneda"`
	MaximoUsuarios  int    `json:"maximo_usuarios"`
	Estado          string `json:"estado"`
}

type adminCredencialesRequest struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
}

type adminCredencialesResponse struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
}

type adminProfileResponse struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Activo  bool   `json:"activo"`
}

type empresaResponse struct {
	ID              int             `json:"id"`
	Nombre          string          `json:"nombre"`
	DocumentoFiscal string          `json:"documento_fiscal,omitempty"`
	Correo          string          `json:"correo,omitempty"`
	Telefono        string          `json:"telefono,omitempty"`
	Direccion       string          `json:"direccion,omitempty"`
	Ciudad          string          `json:"ciudad,omitempty"`
	Pais            string          `json:"pais,omitempty"`
	Moneda          string          `json:"moneda"`
	MonedaInfo      *monedaResponse `json:"moneda_info,omitempty"`
	MaximoUsuarios  int             `json:"maximo_usuarios"`
	Estado          string          `json:"estado"`
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
		Nombre:  adminActual.Nombre,
		Usuario: adminActual.Usuario,
		Activo:  adminActual.Activo,
	})
}

// CrearEmpresa godoc
// @Summary Alta de empresa + usuario principal
// @Description Crea una empresa y su usuario principal en una sola operacion. Requiere autenticacion de administrador.
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
		Nombre:          req.Empresa.Nombre,
		DocumentoFiscal: req.Empresa.DocumentoFiscal,
		Correo:          req.Empresa.Correo,
		Telefono:        req.Empresa.Telefono,
		Direccion:       req.Empresa.Direccion,
		Ciudad:          req.Empresa.Ciudad,
		Pais:            req.Empresa.Pais,
		Moneda:          defaultString(req.Empresa.Moneda, "PEN"),
		MaximoUsuarios:  req.Empresa.MaximoUsuarios,
		Estado:          req.Empresa.Estado,
	}
	u := &domain.Usuario{
		Nombres:        req.Usuario.Nombres,
		Apellidos:      req.Usuario.Apellidos,
		Correo:         req.Usuario.Correo,
		Telefono:       req.Usuario.Telefono,
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
// @Description Devuelve el detalle completo de una empresa, incluyendo la informacion de moneda si esta disponible.
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
	return c.JSON(mapEmpresa(empresa))
}

// ListarEmpresas godoc
// @Summary Listado de empresas
// @Description Retorna todas las empresas registradas para uso del panel administrativo.
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
		mapped := mapEmpresa(e)
		if mapped != nil {
			resp = append(resp, *mapped)
		}
	}
	return c.JSON(resp)
}

// ActualizarEmpresa godoc
// @Summary Actualizar empresa
// @Description Actualiza los datos generales de una empresa existente. Requiere autenticacion de administrador.
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
		ID:              id,
		Nombre:          req.Nombre,
		DocumentoFiscal: req.DocumentoFiscal,
		Correo:          req.Correo,
		Telefono:        req.Telefono,
		Direccion:       req.Direccion,
		Ciudad:          req.Ciudad,
		Pais:            req.Pais,
		Moneda:          defaultString(req.Moneda, "PEN"),
		MaximoUsuarios:  req.MaximoUsuarios,
		Estado:          req.Estado,
	})
	if err != nil {
		if err == service.ErrMonedaInvalida {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(mapEmpresa(empresa))
}

// EliminarEmpresa godoc
// @Summary Eliminar empresa
// @Description Elimina una empresa por ID. Requiere autenticacion de administrador.
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

// ActualizarCredencialesAdmin godoc
// @Summary Cambiar usuario y contraseña del admin autenticado
// @Description Actualiza las credenciales del administrador autenticado.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body adminCredencialesRequest true "Nuevo usuario y contraseña"
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
	adminID, ok := adminIDVal.(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	var req adminCredencialesRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	adminActualizado, err := h.svc.ActualizarCredenciales(c.Context(), adminID, req.Usuario, req.Contrasena)
	if err != nil {
		if err == service.ErrUsuarioAdminInvalido || err == service.ErrContrasenaInvalida {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}
		return fiber.ErrInternalServerError
	}

	return c.JSON(adminCredencialesResponse{
		ID:      adminActualizado.ID,
		Nombre:  adminActualizado.Nombre,
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
	if req.Empresa.Estado != "" && !estadoEmpresaValido(req.Empresa.Estado) {
		return fiber.NewError(http.StatusBadRequest, "empresa.estado debe ser activa, inactiva o suspendida")
	}
	if strings.TrimSpace(req.Usuario.Nombres) == "" {
		return fiber.NewError(http.StatusBadRequest, "usuario.nombres es obligatorio")
	}
	if strings.TrimSpace(req.Usuario.Correo) == "" {
		return fiber.NewError(http.StatusBadRequest, "usuario.correo es obligatorio")
	}
	if strings.TrimSpace(req.Usuario.Password) == "" {
		return fiber.NewError(http.StatusBadRequest, "usuario.password es obligatorio")
	}
	return nil
}

func estadoEmpresaValido(val string) bool {
	switch strings.TrimSpace(val) {
	case "activa", "inactiva", "suspendida":
		return true
	default:
		return false
	}
}

func parseIDParam(c *fiber.Ctx) (int, error) {
	return strconv.Atoi(c.Params("id"))
}
