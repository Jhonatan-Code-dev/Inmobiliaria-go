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

// RegistrarYCobrar godoc
// @Summary Registrar lectura y generar deuda automáticamente
// @Tags Servicios
// @Router /api/user/servicios/registrar-y-cobrar [post]
func (h *ServicioMedicionController) RegistrarYCobrar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	
	// Estructura temporal para ser flexibles con los tipos de entrada
	type registroReq struct {
		ContratoID      interface{} `json:"contrato_id"`
		TipoServicio    string      `json:"tipo_servicio"`
		LecturaAnterior *float64    `json:"lectura_anterior"`
		LecturaActual   float64     `json:"lectura_actual"`
		PrecioUnitario  float64     `json:"precio_unitario"`
		Factor          float64     `json:"factor"`
		CargoFijo       float64     `json:"cargo_fijo"`
		FechaLectura    string      `json:"fecha_lectura"`
	}

	var raw registroReq
	if err := c.BodyParser(&raw); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{
			Message: "Cuerpo de petición inválido. Asegúrese de enviar los campos correctamente.",
		})
	}

	// Convertir ContratoID de forma robusta
	var contratoID int
	switch v := raw.ContratoID.(type) {
	case float64:
		contratoID = int(v)
	case string:
		contratoID, _ = strconv.Atoi(v)
	case int:
		contratoID = v
	default:
		return c.Status(http.StatusBadRequest).JSON(errorResponse{
			Message: "contrato_id debe ser un número.",
		})
	}

	if contratoID == 0 {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{
			Message: "contrato_id es obligatorio y debe ser válido.",
		})
	}

	req := domain.RegistroLectura{
		ContratoID:      contratoID,
		TipoServicio:    raw.TipoServicio,
		LecturaAnterior: raw.LecturaAnterior,
		LecturaActual:   raw.LecturaActual,
		PrecioUnitario:  raw.PrecioUnitario,
		Factor:          raw.Factor,
		CargoFijo:       raw.CargoFijo,
		FechaLectura:    raw.FechaLectura,
	}

	med, err := h.svc.RegistrarYCobrar(c.Context(), &req, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(med)
}

// ObtenerUltima godoc
// @Summary Obtener última lectura registrada para un contrato
// @Tags Servicios
// @Router /api/user/servicios/ultimo/{contrato_id} [get]
func (h *ServicioMedicionController) ObtenerUltima(c *fiber.Ctx) error {
	contratoID, _ := c.ParamsInt("contrato_id")
	tipo := c.Query("tipo", "luz")

	med, err := h.svc.ObtenerUltimaLectura(c.Context(), contratoID, tipo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	if med == nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{"lectura_actual": 0})
	}

	return c.JSON(med)
}

// RegistrarMasivo godoc
// @Summary Registrar múltiples lecturas a la vez
// @Tags Servicios
// @Router /api/user/servicios/masivo [post]
func (h *ServicioMedicionController) RegistrarMasivo(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)

	type registroReq struct {
		ContratoID      interface{} `json:"contrato_id"`
		TipoServicio    string      `json:"tipo_servicio"`
		LecturaAnterior *float64    `json:"lectura_anterior"`
		LecturaActual   float64     `json:"lectura_actual"`
		PrecioUnitario  float64     `json:"precio_unitario"`
		Factor          float64     `json:"factor"`
		CargoFijo       float64     `json:"cargo_fijo"`
		FechaLectura    string      `json:"fecha_lectura"`
	}

	var rawItems []registroReq
	if err := c.BodyParser(&rawItems); err != nil {
		return c.Status(http.StatusBadRequest).JSON(errorResponse{
			Message: "Formato masivo inválido. Se espera un array de objetos.",
		})
	}

	var registros []domain.RegistroLectura
	for _, raw := range rawItems {
		var contratoID int
		switch v := raw.ContratoID.(type) {
		case float64:
			contratoID = int(v)
		case string:
			contratoID, _ = strconv.Atoi(v)
		case int:
			contratoID = v
		}

		if contratoID > 0 {
			registros = append(registros, domain.RegistroLectura{
				ContratoID:      contratoID,
				TipoServicio:    raw.TipoServicio,
				LecturaAnterior: raw.LecturaAnterior,
				LecturaActual:   raw.LecturaActual,
				PrecioUnitario:  raw.PrecioUnitario,
				Factor:          raw.Factor,
				CargoFijo:       raw.CargoFijo,
				FechaLectura:    raw.FechaLectura,
			})
		}
	}

	mediciones, err := h.svc.RegistrarMasivo(c.Context(), registros, empresaID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(mediciones)
}

// ListarPendientes godoc
// @Summary Listar contratos que faltan registrar lectura este mes
// @Tags Servicios
// @Router /api/user/servicios/pendientes [get]
func (h *ServicioMedicionController) ListarPendientes(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	tipo := c.Query("tipo", "luz")

	pendientes, err := h.svc.ListarPendientesLectura(c.Context(), empresaID, tipo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(pendientes)
}
