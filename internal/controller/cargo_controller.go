package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CargoController struct {
	svc domain.CargoService
}

func NewCargoController(svc domain.CargoService) *CargoController {
	return &CargoController{svc: svc}
}

// Listar godoc
// @Summary Listar cargos de la empresa
// @Tags Cargos
// @Router /api/user/cargos [get]
func (h *CargoController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "1"))
	porPagina, _ := strconv.Atoi(c.Query("por_pagina", "10"))
	contratoID, _ := strconv.Atoi(c.Query("contrato_id", "0"))
	estado := c.Query("estado", "")

	filtros := domain.CargoFiltros{
		EmpresaID:  empresaID,
		ContratoID: contratoID,
		Estado:     estado,
		Pagina:     pag,
		PorPagina:  porPagina,
	}

	cargos, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(paginatedResponse{
		Datos: cargos,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    porPagina,
			Paginas:      (total + porPagina - 1) / porPagina,
		},
	})
}

// Obtener godoc
// @Summary Obtener detalle de un cargo
// @Tags Cargos
// @Router /api/user/cargos/{id} [get]
func (h *CargoController) Obtener(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	cargo, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(errorResponse{Message: "cargo no encontrado"})
	}

	return c.JSON(cargo)
}

// Crear godoc
// @Summary Registrar nuevo cargo manual
// @Tags Cargos
// @Router /api/user/cargos [post]
func (h *CargoController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.RegistroCargo
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	cargo, err := h.svc.Crear(c.Context(), &req, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(cargo)
}

// Actualizar godoc
// @Summary Editar cargo
// @Tags Cargos
// @Router /api/user/cargos/{id} [put]
func (h *CargoController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	var req domain.RegistroCargo
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	cargo, err := h.svc.Actualizar(c.Context(), id, empresaID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(cargo)
}

// Eliminar godoc
// @Summary Eliminar cargo
// @Tags Cargos
// @Router /api/user/cargos/{id} [delete]
func (h *CargoController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "cargo eliminado"})
}
