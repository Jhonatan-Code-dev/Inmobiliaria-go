package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ServicioMedicionController struct {
	svc domain.ServicioMedicionService
}

func NewServicioMedicionController(svc domain.ServicioMedicionService) *ServicioMedicionController {
	return &ServicioMedicionController{svc: svc}
}

// Listar godoc
// @Summary Listar mediciones de servicios
// @Tags Servicios
// @Router /api/user/servicios [get]
func (h *ServicioMedicionController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "1"))
	porPagina, _ := strconv.Atoi(c.Query("por_pagina", "10"))
	contratoID, _ := strconv.Atoi(c.Query("contrato_id", "0"))

	filtros := domain.ServicioMedicionFiltros{
		EmpresaID:  empresaID,
		ContratoID: contratoID,
		Pagina:     pag,
		PorPagina:  porPagina,
	}

	mediciones, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(paginatedResponse{
		Datos: mediciones,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    porPagina,
			Paginas:      (total + porPagina - 1) / porPagina,
		},
	})
}

// Obtener godoc
// @Summary Obtener detalle de medición
// @Tags Servicios
// @Router /api/user/servicios/{id} [get]
func (h *ServicioMedicionController) Obtener(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	med, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(errorResponse{Message: "medición no encontrada"})
	}

	return c.JSON(med)
}

// Crear godoc
// @Summary Registrar lectura de servicio
// @Tags Servicios
// @Router /api/user/servicios [post]
func (h *ServicioMedicionController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.RegistroLectura
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	med, err := h.svc.Registrar(c.Context(), &req, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(med)
}

// Eliminar godoc
// @Summary Eliminar medición
// @Tags Servicios
// @Router /api/user/servicios/{id} [delete]
func (h *ServicioMedicionController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "medición eliminada"})
}

// Actualizar godoc
// @Summary Editar lectura errónea
// @Tags Servicios
// @Param id path int true "ID de la medición"
// @Param request body map[string]float64 true "Nueva lectura_actual"
// @Router /api/user/servicios/{id} [put]
func (h *ServicioMedicionController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	
	type updateReq struct {
		LecturaActual float64 `json:"lectura_actual"`
	}
	var req updateReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	updated, err := h.svc.Actualizar(c.Context(), id, empresaID, req.LecturaActual)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(updated)
}
