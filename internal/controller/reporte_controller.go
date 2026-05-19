package controller

import (
	"time"

	"rentals-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// ReporteController expone los endpoints de reportes gerenciales para gráficos.
type ReporteController struct {
	svc domain.ReporteService
}

func NewReporteController(svc domain.ReporteService) *ReporteController {
	return &ReporteController{svc: svc}
}

// helper to parse dates from fiber query parameters
func parseRangoFechasQuery(c *fiber.Ctx) (time.Time, time.Time, error) {
	var desde, hasta time.Time

	if desdeStr := c.Query("desde"); desdeStr != "" {
		t, err := time.Parse("2006-01-02", desdeStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		desde = t.UTC()
	}

	if hastaStr := c.Query("hasta"); hastaStr != "" {
		t, err := time.Parse("2006-01-02", hastaStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		hasta = t.UTC()
	}

	return desde, hasta, nil
}

// IngresosGastos godoc
// @Summary      Reporte de Ingresos vs Gastos
// @Description  Retorna la serie mensual del balance neto de ingresos y gastos en el rango de fechas indicado para gráficos de barras o líneas. Si no se envían fechas, usa los últimos 12 meses.
// @Tags         Reportes Gerenciales
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {object}  domain.ReporteIngresosGastos
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/reportes/ingresos-gastos [get]
func (h *ReporteController) IngresosGastos(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	desde, hasta, err := parseRangoFechasQuery(c)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de fecha inválido, usa YYYY-MM-DD"})
	}

	data, err := h.svc.ObtenerIngresosGastos(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener reporte de ingresos y gastos: " + err.Error()})
	}

	return c.JSON(data)
}

// MetodosPago godoc
// @Summary      Distribución de ingresos por método de pago
// @Description  Detalla el volumen de ingresos agrupado por método de pago para gráficos circulares/torta (Pie / Doughnut). Si no se envían fechas, usa los últimos 12 meses.
// @Tags         Reportes Gerenciales
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {array}   domain.DistribucionMetodoPago
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/reportes/metodos-pago [get]
func (h *ReporteController) MetodosPago(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	desde, hasta, err := parseRangoFechasQuery(c)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de fecha inválido, usa YYYY-MM-DD"})
	}

	data, err := h.svc.ObtenerDistribucionMetodosPago(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener distribución por método de pago: " + err.Error()})
	}

	return c.JSON(data)
}

// CategoriasGastos godoc
// @Summary      Distribución de gastos por categoría
// @Description  Detalla el volumen de egresos agrupado por la categoría o tipo de pago del gasto para gráficos circulares/torta (Pie / Doughnut). Si no se envían fechas, usa los últimos 12 meses.
// @Tags         Reportes Gerenciales
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {array}   domain.DistribucionCategoriaGasto
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/reportes/categorias-gastos [get]
func (h *ReporteController) CategoriasGastos(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	desde, hasta, err := parseRangoFechasQuery(c)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de fecha inválido, usa YYYY-MM-DD"})
	}

	data, err := h.svc.ObtenerDistribucionCategoriasGastos(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener distribución por categoría: " + err.Error()})
	}

	return c.JSON(data)
}

// RentabilidadPropiedades godoc
// @Summary      Rentabilidad y ocupación por propiedad
// @Description  Devuelve un análisis comparativo por propiedad que incluye tasa de ocupación, ingresos generados, gastos prorrateados asignados y rentabilidad neta para gráficos de barras horizontales o comparativos. Si no se envían fechas, usa los últimos 12 meses.
// @Tags         Reportes Gerenciales
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {array}   domain.RentabilidadPropiedad
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/reportes/rentabilidad-propiedades [get]
func (h *ReporteController) RentabilidadPropiedades(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	desde, hasta, err := parseRangoFechasQuery(c)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de fecha inválido, usa YYYY-MM-DD"})
	}

	data, err := h.svc.ObtenerRentabilidadPropiedades(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener rentabilidad de propiedades: " + err.Error()})
	}

	return c.JSON(data)
}

// TicketsMantenimiento godoc
// @Summary      Resumen de tickets de soporte y mantenimiento
// @Description  Detalla métricas del módulo de soporte/mantenimiento agrupados por estado y prioridad creados en el rango indicado. Útil para gráficos de barras apiladas. Si no se envían fechas, usa los últimos 12 meses.
// @Tags         Reportes Gerenciales
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {object}  domain.ResumenMantenimientoReporte
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/reportes/tickets-mantenimiento [get]
func (h *ReporteController) TicketsMantenimiento(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	desde, hasta, err := parseRangoFechasQuery(c)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de fecha inválido, usa YYYY-MM-DD"})
	}

	data, err := h.svc.ObtenerResumenMantenimiento(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener resumen de mantenimiento: " + err.Error()})
	}

	return c.JSON(data)
}
