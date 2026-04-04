package controller

import (
	"net/http"
	"strconv"

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

// CrearEmpresa godoc
// @Summary Alta de empresa + usuario principal
// @Tags admin
// @Accept json
// @Produce json
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
	if req.Empresa.Nombre == "" || req.Usuario.Correo == "" || req.Usuario.Password == "" {
		return fiber.ErrBadRequest
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
// @Tags admin
// @Produce json
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
// @Tags admin
// @Produce json
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
// @Tags admin
// @Accept json
// @Produce json
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
// @Tags admin
// @Produce json
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
// @Tags admin
// @Accept json
// @Produce json
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

func parseIDParam(c *fiber.Ctx) (int, error) {
	return strconv.Atoi(c.Params("id"))
}
