package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CitasController struct {
	svc         domain.CitaService
	inmuebleSvc domain.InmuebleService
	clienteSvc  domain.ClienteService
}

func NewCitasController(svc domain.CitaService, inmuebleSvc domain.InmuebleService, clienteSvc domain.ClienteService) *CitasController {
	return &CitasController{
		svc:         svc,
		inmuebleSvc: inmuebleSvc,
		clienteSvc:  clienteSvc,
	}
}

// Listar godoc
// @Summary Listar citas / visitas de la empresa
// @Tags Citas
// @Router /api/user/citas [get]
func (h *CitasController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "0"))
	porPagina, _ := strconv.Atoi(c.Query("por_pagina", "0"))
	propiedadID, _ := strconv.Atoi(c.Query("propiedad_id", "0"))
	unidadID, _ := strconv.Atoi(c.Query("unidad_id", "0"))
	estado := c.Query("estado")

	var desde *time.Time
	if d := c.Query("desde"); d != "" {
		if t, err := time.Parse(time.RFC3339, d); err == nil {
			desde = &t
		} else if t, err := time.Parse("2006-01-02", d); err == nil {
			desde = &t
		}
	}
	var hasta *time.Time
	if hStr := c.Query("hasta"); hStr != "" {
		if t, err := time.Parse(time.RFC3339, hStr); err == nil {
			hasta = &t
		} else if t, err := time.Parse("2006-01-02", hStr); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			hasta = &t
		}
	}

	filtros := domain.CitaFiltros{
		EmpresaID:   empresaID,
		PropiedadID: propiedadID,
		UnidadID:    unidadID,
		Estado:      estado,
		Busqueda:    c.Query("buscar"),
		Desde:       desde,
		Hasta:       hasta,
		Pagina:      pag,
		PorPagina:   porPagina,
	}

	citas, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	// Si no se paginó, retornamos los datos directamente en un map simple
	if pag <= 0 || porPagina <= 0 {
		return c.JSON(fiber.Map{
			"datos": citas,
			"total": total,
		})
	}

	return c.JSON(paginatedResponse{
		Datos: citas,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    porPagina,
			Paginas:      (total + porPagina - 1) / porPagina,
		},
	})
}

// Obtener godoc
// @Summary Obtener detalle de cita
// @Tags Citas
// @Router /api/user/citas/{id} [get]
func (h *CitasController) Obtener(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	citaObj, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(errorResponse{Message: "cita no encontrada"})
	}

	return c.JSON(citaObj)
}

// Crear godoc
// @Summary Agendar nueva cita / visita
// @Tags Citas
// @Router /api/user/citas [post]
func (h *CitasController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.RegistroCita
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	citaObj, err := h.svc.Crear(c.Context(), &req, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(citaObj)
}

// Actualizar godoc
// @Summary Actualizar cita
// @Tags Citas
// @Router /api/user/citas/{id} [put]
func (h *CitasController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	var req domain.RegistroCita
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	citaObj, err := h.svc.Actualizar(c.Context(), id, empresaID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(citaObj)
}

// CambiarEstado godoc
// @Summary Cambiar estado de cita
// @Tags Citas
// @Router /api/user/citas/{id}/estado [patch]
func (h *CitasController) CambiarEstado(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	type EstadoReq struct {
		Estado string `json:"estado"`
	}
	var req EstadoReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Estado == "" {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: "el estado es requerido"})
	}

	citaObj, err := h.svc.CambiarEstado(c.Context(), id, empresaID, req.Estado)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    citaObj,
		"message": "Estado de la cita actualizado",
	})
}

// Eliminar godoc
// @Summary Eliminar cita
// @Tags Citas
// @Router /api/user/citas/{id} [delete]
func (h *CitasController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "cita eliminada"})
}

// ConfigFormulario godoc
// @Summary Catálogos para el formulario de citas
// @Tags Citas
// @Router /api/user/citas/config-formulario [get]
func (h *CitasController) ConfigFormulario(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)

	// 1. Obtener Inmuebles
	inmuebles, _, err := h.inmuebleSvc.Listar(c.Context(), domain.InmuebleFiltros{
		EmpresaID: empresaID,
		Limite:    100,
	})
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "Error al obtener inmuebles: " + err.Error()})
	}

	listaInmuebles := make([]fiber.Map, 0, len(inmuebles))
	for _, in := range inmuebles {
		listaInmuebles = append(listaInmuebles, fiber.Map{
			"id":     in.ID,
			"nombre": in.Nombre,
		})
	}

	// 2. Obtener Clientes
	clientes, _, err := h.clienteSvc.Listar(c.Context(), domain.ClienteFiltros{
		EmpresaID: empresaID,
		Limite:    200,
	})
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "Error al obtener clientes: " + err.Error()})
	}

	listaClientes := make([]fiber.Map, 0, len(clientes))
	for _, cl := range clientes {
		nombreCompleto := cl.Nombres
		if cl.Apellidos != nil {
			nombreCompleto += " " + *cl.Apellidos
		}
		listaClientes = append(listaClientes, fiber.Map{
			"id":     cl.ID,
			"nombre": nombreCompleto,
		})
	}

	return c.JSON(fiber.Map{
		"inmuebles": listaInmuebles,
		"clientes":  listaClientes,
		"estados":   []string{"programada", "realizada", "cancelada", "no_asistio"},
	})
}
