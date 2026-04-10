package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TicketController struct {
	svc domain.TicketService
}

func NewTicketController(svc domain.TicketService) *TicketController {
	return &TicketController{svc: svc}
}

// Listar godoc
// @Summary Listar tickets de mantenimiento
// @Tags Tickets
// @Router /api/user/tickets [get]
func (h *TicketController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "1"))
	porPagina, _ := strconv.Atoi(c.Query("por_pagina", "10"))
	unidadID, _ := strconv.Atoi(c.Query("unidad_id", "0"))
	estado := c.Query("estado")

	filtros := domain.TicketFiltros{
		EmpresaID:  empresaID,
		UnidadID:   unidadID,
		Estado:     estado,
		Pagina:     pag,
		PorPagina:  porPagina,
	}

	tickets, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(paginatedResponse{
		Datos: tickets,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    porPagina,
			Paginas:      (total + porPagina - 1) / porPagina,
		},
	})
}

// Obtener godoc
// @Summary Obtener detalle de ticket
// @Tags Tickets
// @Router /api/user/tickets/{id} [get]
func (h *TicketController) Obtener(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	t, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(errorResponse{Message: "ticket no encontrado"})
	}

	return c.JSON(t)
}

// Crear godoc
// @Summary Crear ticket de mantenimiento
// @Tags Tickets
// @Router /api/user/tickets [post]
func (h *TicketController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.RegistroTicket
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	t, err := h.svc.Crear(c.Context(), &req, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(t)
}

// Actualizar godoc
// @Summary Actualizar ticket
// @Tags Tickets
// @Router /api/user/tickets/{id} [put]
func (h *TicketController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	
	type UpdateReq struct {
		domain.RegistroTicket
		Estado string `json:"estado"`
	}
	var req UpdateReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	t, err := h.svc.Actualizar(c.Context(), id, empresaID, &req.RegistroTicket, req.Estado)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(t)
}

// Eliminar godoc
// @Summary Eliminar ticket
// @Tags Tickets
// @Router /api/user/tickets/{id} [delete]
func (h *TicketController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "ticket eliminado"})
}
