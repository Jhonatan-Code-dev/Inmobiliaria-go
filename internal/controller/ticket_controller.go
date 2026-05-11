package controller

import (
	"net/http"
	"rentals-go/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TicketController struct {
	svc         domain.TicketService
	inmuebleSvc domain.InmuebleService
	clienteSvc  domain.ClienteService
}

func NewTicketController(svc domain.TicketService, inmuebleSvc domain.InmuebleService, clienteSvc domain.ClienteService) *TicketController {
	return &TicketController{
		svc:         svc,
		inmuebleSvc: inmuebleSvc,
		clienteSvc:  clienteSvc,
	}
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
	propiedadID, _ := strconv.Atoi(c.Query("propiedad_id", "0"))
	estado := c.Query("estado")

	filtros := domain.TicketFiltros{
		EmpresaID:   empresaID,
		PropiedadID: propiedadID,
		UnidadID:    unidadID,
		Estado:     estado,
		Busqueda:   c.Query("buscar"),
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

// Resumen godoc
// @Summary Obtener resumen de tickets
// @Description Retorna los totales de tickets agrupados por estado
// @Tags Tickets
// @Security BearerAuth
// @Param propiedad_id query int false "ID del Inmueble/Propiedad"
// @Success 200 {object} domain.TicketResumen
// @Router /api/user/tickets/resumen [get]
func (h *TicketController) Resumen(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	propiedadID, _ := strconv.Atoi(c.Query("propiedad_id", "0"))

	resumen, err := h.svc.ObtenerResumen(c.Context(), empresaID, propiedadID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(resumen)
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
// @Security BearerAuth
// @Router /api/user/tickets/{id} [delete]
func (h *TicketController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "ticket eliminado"})
}

// CambiarEstado godoc
// @Summary Cambiar estado de ticket
// @Description Cambia el estado de un ticket a en_progreso, resuelto, o cerrado
// @Tags Tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del Ticket"
// @Param req body domain.CambiarEstadoTicket true "Nuevo estado"
// @Success 200 {object} domain.Ticket
// @Router /api/user/tickets/{id}/estado [patch]
func (h *TicketController) CambiarEstado(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	var req domain.CambiarEstadoTicket
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: "Cuerpo de la petición inválido"})
	}

	if req.Estado == "" {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{Message: "estado es requerido"})
	}

	t, err := h.svc.CambiarEstado(c.Context(), id, empresaID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    t,
		"message": "Estado del ticket actualizado exitosamente",
	})
}

// ConfigFormulario godoc
// @Summary Obtener catálogos para el formulario de tickets
// @Description Retorna las listas de inmuebles, clientes, prioridades y estados para llenar los dropdowns del formulario.
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/user/tickets/config-formulario [get]
func (h *TicketController) ConfigFormulario(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)

	// 1. Obtener Inmuebles
	inmuebles, _, err := h.inmuebleSvc.Listar(c.Context(), domain.InmuebleFiltros{
		EmpresaID: empresaID,
		Limite:    100, // Un número razonable para un dropdown
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
		"inmuebles":   listaInmuebles,
		"clientes":    listaClientes,
		"prioridades": []string{"baja", "media", "alta"},
		"estados":     []string{"abierto", "en_progreso", "resuelto", "cerrado"},
	})
}

// ColaTrabajo godoc
// @Summary Listar cola de trabajo (Pendientes prioritarios)
// @Description Obtiene los tickets abiertos o en progreso, ordenados por prioridad (Alta -> Baja) y antigüedad.
// @Tags Tickets
// @Security BearerAuth
// @Produce json
// @Param pag query int false "Página" default(1)
// @Param buscar query string false "Búsqueda"
// @Param propiedad_id query int false "Filtrar por inmueble"
// @Success 200 {object} paginatedResponse
// @Router /api/user/tickets/cola-trabajo [get]
func (h *TicketController) ColaTrabajo(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	pag, _ := strconv.Atoi(c.Query("pag", "1"))
	propiedadID, _ := strconv.Atoi(c.Query("propiedad_id", "0"))

	filtros := domain.TicketFiltros{
		EmpresaID:   empresaID,
		PropiedadID: propiedadID,
		Busqueda:    c.Query("buscar"),
		Pagina:      pag,
		PorPagina:   10,
	}

	tickets, total, err := h.svc.ListarColaTrabajo(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(paginatedResponse{
		Datos: tickets,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: pag,
			PorPagina:    filtros.PorPagina,
			Paginas:      (total + filtros.PorPagina - 1) / filtros.PorPagina,
		},
	})
}

