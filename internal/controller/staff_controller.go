package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type StaffController struct {
	svc domain.StaffService
}

func NewStaffController(svc domain.StaffService) *StaffController {
	return &StaffController{svc: svc}
}

// Listar godoc
// @Summary Listar staff de la empresa
// @Tags Staff
// @Param empresa_id query int true "ID de la empresa"
// @Param pag query int false "Página"
// @Param por_pagina query int false "Items por página"
// @Param buscar query string false "Buscar por nombre de usuario"
// @Router /api/user/staff [get]
func (h *StaffController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "1"))
	porPagina, _ := strconv.Atoi(c.Query("por_pagina", "10"))
	buscar := c.Query("buscar", "")

	filtros := domain.StaffFiltros{
		EmpresaID: empresaID,
		Pagina:    pag,
		PorPagina: porPagina,
		Busqueda:  buscar,
	}

	staff, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(paginatedResponse{
		Datos: staff,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    porPagina,
			Paginas:      (total + porPagina - 1) / porPagina,
		},
	})
}

// Obtener godoc
// @Summary Obtener detalle de un empleado
// @Tags Staff
// @Router /api/user/staff/{id} [get]
func (h *StaffController) Obtener(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	staff, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(errorResponse{Message: "empleado no encontrado"})
	}

	return c.JSON(staff)
}

// Crear godoc
// @Summary Registrar nuevo empleado
// @Tags Staff
// @Router /api/user/staff [post]
func (h *StaffController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.RegistroStaff
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	req.EmpresaID = empresaID

	staff, err := h.svc.Registrar(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(staff)
}

type actualizarStaffRequest struct {
	RolID  int    `json:"rol_id"`
	Estado string `json:"estado"`
}

// Actualizar godoc
// @Summary Editar empleado
// @Tags Staff
// @Router /api/user/staff/{id} [put]
func (h *StaffController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	var req actualizarStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	staff, err := h.svc.Actualizar(c.Context(), id, empresaID, req.RolID, req.Estado)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(staff)
}

// Eliminar godoc
// @Summary Eliminar/Dar de baja empleado
// @Tags Staff
// @Router /api/user/staff/{id} [delete]
func (h *StaffController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "empleado eliminado"})
}

// ListarRoles godoc
// @Summary Listar roles disponibles
// @Tags Staff
// @Router /api/user/staff/roles [get]
func (h *StaffController) ListarRoles(c *fiber.Ctx) error {
	roles, err := h.svc.ListarRoles(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(roles)
}
