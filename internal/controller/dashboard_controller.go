package controller

import (
	"time"

	"rentals-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// DashboardController expone todos los endpoints de KPIs / Dashboard.
type DashboardController struct {
	svc domain.DashboardService
}

func NewDashboardController(svc domain.DashboardService) *DashboardController {
	return &DashboardController{svc: svc}
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard
// ─────────────────────────────────────────────

// ResumenGeneral godoc
// @Summary      KPIs generales del negocio
// @Description  Retorna los indicadores clave del mes actual: ocupación, ingresos, gastos, morosidad y contratos activos. Es el endpoint principal para el panel de control (dashboard).
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int  true  "ID de la empresa"
// @Success      200  {object}  domain.ResumenGeneral
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard [get]
func (h *DashboardController) ResumenGeneral(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	data, err := h.svc.ObtenerResumenGeneral(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener resumen: " + err.Error()})
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/ocupacion
// ─────────────────────────────────────────────

// Ocupacion godoc
// @Summary      Tasa de ocupación por propiedad
// @Description  Devuelve la tasa de ocupación global y el desglose por cada propiedad: unidades totales, ocupadas y libres.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int  true  "ID de la empresa"
// @Success      200  {object}  domain.ResumenOcupacion
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/ocupacion [get]
func (h *DashboardController) Ocupacion(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	data, err := h.svc.ObtenerOcupacion(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener ocupación: " + err.Error()})
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/morosidad
// ─────────────────────────────────────────────

// Morosidad godoc
// @Summary      Reporte de inquilinos morosos
// @Description  Lista todos los inquilinos con cargos vencidos sin pagar, con el monto adeudado y los días de atraso.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int  true  "ID de la empresa"
// @Success      200  {object}  domain.ResumenMorosidad
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/morosidad [get]
func (h *DashboardController) Morosidad(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	data, err := h.svc.ObtenerMorosidad(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener morosidad: " + err.Error()})
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/financiero
// ─────────────────────────────────────────────

// ReporteFinanciero godoc
// @Summary      Reporte financiero por rango de fechas
// @Description  Devuelve el total de ingresos, gastos y balance neto en el rango indicado, además de la serie mensual desglosada. Si no se envían fechas, usa los últimos 6 meses.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Success      200  {object}  domain.ReporteFinanciero
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/financiero [get]
func (h *DashboardController) ReporteFinanciero(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	var desde, hasta time.Time
	if desdeStr := c.Query("desde"); desdeStr != "" {
		t, err := time.Parse("2006-01-02", desdeStr)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "desde: formato inválido, usa YYYY-MM-DD"})
		}
		desde = t.UTC()
	}
	if hastaStr := c.Query("hasta"); hastaStr != "" {
		t, err := time.Parse("2006-01-02", hastaStr)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "hasta: formato inválido, usa YYYY-MM-DD"})
		}
		hasta = t.UTC()
	}

	data, err := h.svc.ObtenerReporteFinanciero(c.Context(), empresaID, desde, hasta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error en reporte financiero: " + err.Error()})
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/contratos-por-vencer
// ─────────────────────────────────────────────

// ContratosProximosVencer godoc
// @Summary      Contratos próximos a vencer
// @Description  Lista los contratos activos cuya fecha_fin es menor o igual a los próximos N días. Por defecto 30 días.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int  true   "ID de la empresa"
// @Param        dias        query  int  false  "Días de anticipación (default: 30)"
// @Success      200  {array}   domain.ContratoProximoVencer
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/contratos-por-vencer [get]
func (h *DashboardController) ContratosProximosVencer(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	dias := c.QueryInt("dias", 30)
	data, err := h.svc.ObtenerContratosProximosVencer(c.Context(), empresaID, dias)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener contratos: " + err.Error()})
	}
	if data == nil {
		data = []domain.ContratoProximoVencer{}
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/estado-cuenta/:clienteId
// ─────────────────────────────────────────────

// EstadoCuentaCliente godoc
// @Summary      Estado de cuenta de un inquilino
// @Description  Devuelve el resumen financiero completo del cliente: cargos emitidos, monto pagado y saldo pendiente, con el desglose cargo por cargo.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        clienteId   path   int  true  "ID del cliente (inquilino)"
// @Param        empresa_id  query  int  true  "ID de la empresa"
// @Success      200  {object}  domain.EstadoCuentaCliente
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/estado-cuenta/{clienteId} [get]
func (h *DashboardController) EstadoCuentaCliente(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	clienteID, err := c.ParamsInt("clienteId")
	if err != nil || clienteID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "clienteId inválido"})
	}
	data, err := h.svc.ObtenerEstadoCuentaCliente(c.Context(), empresaID, clienteID)
	if err != nil {
		return manejarErrorCliente(err, "cliente")
	}
	return c.JSON(data)
}

// ─────────────────────────────────────────────
// GET /api/user/dashboard/top-unidades
// ─────────────────────────────────────────────

// TopUnidades godoc
// @Summary      Top unidades por ingresos generados
// @Description  Devuelve las unidades con mayor recaudación en el rango de fechas indicado. Por defecto usa el mes actual y muestra las 10 primeras.
// @Tags         Dashboard
// @Security     BearerAuth
// @Produce      json
// @Param        empresa_id  query  int     true   "ID de la empresa"
// @Param        desde       query  string  false  "Fecha inicio (YYYY-MM-DD)"
// @Param        hasta       query  string  false  "Fecha fin (YYYY-MM-DD)"
// @Param        limite      query  int     false  "Número de resultados (default: 10)"
// @Success      200  {array}   domain.TopUnidad
// @Failure      400  {object}  errorResponse
// @Failure      401  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/user/dashboard/top-unidades [get]
func (h *DashboardController) TopUnidades(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	var desde, hasta time.Time
	if desdeStr := c.Query("desde"); desdeStr != "" {
		t, err := time.Parse("2006-01-02", desdeStr)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "desde: formato inválido, usa YYYY-MM-DD"})
		}
		desde = t.UTC()
	}
	if hastaStr := c.Query("hasta"); hastaStr != "" {
		t, err := time.Parse("2006-01-02", hastaStr)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "hasta: formato inválido, usa YYYY-MM-DD"})
		}
		hasta = t.UTC()
	}

	limite := c.QueryInt("limite", 10)
	data, err := h.svc.ObtenerTopUnidades(c.Context(), empresaID, desde, hasta, limite)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: "error al obtener top unidades: " + err.Error()})
	}
	if data == nil {
		data = []domain.TopUnidad{}
	}
	return c.JSON(data)
}
