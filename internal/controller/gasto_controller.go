package controller

import (
	"strconv"
	"time"

	"rentals-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type GastoController struct {
	svc domain.GastoService
}

func NewGastoController(svc domain.GastoService) *GastoController {
	return &GastoController{svc: svc}
}

type gastoRequest struct {
	EmpresaID   int     `json:"empresa_id" example:"1"`
	Monto       float64 `json:"monto" example:"50.00"`
	Fecha       string  `json:"fecha" example:"2024-04-07"`
	TipoPagoID  int    `json:"tipo_pago_id" example:"3"`
	Descripcion string `json:"descripcion" example:"Pago de servicios básicos"`
}

type gastoResponse struct {
	ID          int       `json:"id" example:"10"`
	EmpresaID   int       `json:"empresa_id" example:"1"`
	Monto       float64   `json:"monto" example:"50.00"`
	Fecha       time.Time `json:"fecha" example:"2024-04-07T00:00:00Z"`
	TipoPagoID  int       `json:"tipo_pago_id" example:"3"`
	Descripcion string    `json:"descripcion" example:"Pago de servicios básicos"`
}

type tipoPagoResponse struct {
	ID     int    `json:"id" example:"1"`
	Nombre string `json:"nombre" example:"efectivo"`
}

type listadoGastosResponse struct {
	Datos      []gastoResponse   `json:"datos"`
	Paginacion paginadorResponse `json:"paginacion"`
}

// ListarTiposPago godoc
// @Summary Listar los tipos de métodos de pago
// @Description Obtiene la lista de tipos de pago disponibles (Yape, Plin, Efectivo, etc.)
// @Tags Gastos
// @Security BearerAuth
// @Produce json
// @Success 200 {array} tipoPagoResponse
// @Router /api/user/gastos/tipos-pago [get]
func (h *GastoController) ListarTiposPago(c *fiber.Ctx) error {
	list, err := h.svc.ListarTiposPago(c.Context())
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	
	res := make([]tipoPagoResponse, 0, len(list))
	for _, tp := range list {
		res = append(res, tipoPagoResponse{
			ID:     tp.ID,
			Nombre: tp.Nombre,
		})
	}

	return c.JSON(res)
}

// Listar godoc
// @Summary Listar gastos con filtros y paginación
// @Description Obtiene una lista paginada de gastos (máx 10). Filtros opcionales: anio, mes, desde, hasta, fecha.
// @Tags Gastos
// @Security BearerAuth
// @Produce json
// @Param pag query int false "Número de página" default(1)
// @Param anio query int false "Año del gasto"
// @Param mes query int false "Mes del gasto (1-12)"
// @Param desde query string false "Fecha inicio (YYYY-MM-DD)"
// @Param hasta query string false "Fecha fin (YYYY-MM-DD)"
// @Param fecha query string false "Fecha exacta (YYYY-MM-DD)"
// @Success 200 {object} listadoGastosResponse
// @Router /api/user/gastos [get]
func (h *GastoController) Listar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)

	filtros := domain.GastoFiltros{
		EmpresaID: empresaID,
		Pagina:    c.QueryInt("pag", 1),
		Limite:    10,
		Anio:      c.QueryInt("anio"),
		Mes:       c.QueryInt("mes"),
	}

	if f := c.Query("fecha"); f != "" {
		t, _ := time.Parse("2006-01-02", f)
		if !t.IsZero() {
			filtros.Fecha = &t
		}
	}
	if d := c.Query("desde"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		if !t.IsZero() {
			filtros.Desde = &t
		}
	}
	if h_query := c.Query("hasta"); h_query != "" {
		t, _ := time.Parse("2006-01-02", h_query)
		if !t.IsZero() {
			filtros.Hasta = &t
		}
	}

	list, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	datos := make([]gastoResponse, 0, len(list))
	for _, g := range list {
		datos = append(datos, mapGastoToResponse(g))
	}

	paginas := (total + filtros.Limite - 1) / filtros.Limite

	return c.JSON(listadoGastosResponse{
		Datos: datos,
		Paginacion: paginadorResponse{
			Total:     total,
			Paginas:   paginas,
			Pagina:    filtros.Pagina,
			PorPagina: filtros.Limite,
		},
	})
}

// Crear godoc
// @Summary Registrar un nuevo gasto
// @Description Crea un nuevo gasto y un movimiento de egreso en caja automático.
// @Tags Gastos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body gastoRequest true "Datos del gasto"
// @Success 201 {object} gastoResponse
// @Router /api/user/gastos [post]
func (h *GastoController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req gastoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	fecha, _ := time.Parse(time.RFC3339, req.Fecha)
	if fecha.IsZero() {
		fecha, _ = time.Parse("2006-01-02", req.Fecha)
	}

	// Priorizar ID del body si viene, sino usar el de la sesión (Local)
	finalEmpresaID := empresaID
	if req.EmpresaID > 0 {
		finalEmpresaID = req.EmpresaID
	}

	gasto := &domain.Gasto{
		EmpresaID:   finalEmpresaID,
		Monto:       req.Monto,
		Fecha:       fecha,
		TipoPagoID:  req.TipoPagoID,
		Descripcion: req.Descripcion,
	}

	nuevo, err := h.svc.RegistrarGasto(c.Context(), gasto)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(mapGastoToResponse(nuevo))
}

// Actualizar godoc
// @Summary Actualizar un gasto
// @Description Actualiza los datos de un gasto existente.
// @Tags Gastos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del gasto"
// @Param request body gastoRequest true "Datos a actualizar"
// @Success 200 {object} gastoResponse
// @Router /api/user/gastos/{id} [put]
func (h *GastoController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	var req gastoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	fecha, _ := time.Parse(time.RFC3339, req.Fecha)
	if fecha.IsZero() {
		fecha, _ = time.Parse("2006-01-02", req.Fecha)
	}

	gasto := &domain.Gasto{
		ID:          id,
		EmpresaID:   empresaID,
		Monto:       req.Monto,
		Fecha:       fecha,
		TipoPagoID:  req.TipoPagoID,
		Descripcion: req.Descripcion,
	}

	actualizado, err := h.svc.ActualizarGasto(c.Context(), gasto)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(mapGastoToResponse(actualizado))
}

// Eliminar godoc
// @Summary Eliminar un gasto
// @Description Elimina el registro de gasto permanentemente.
// @Tags Gastos
// @Security BearerAuth
// @Param id path int true "ID del gasto"
// @Success 200 {object} map[string]string
// @Router /api/user/gastos/{id} [delete]
func (h *GastoController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	if err := h.svc.EliminarGasto(c.Context(), id, empresaID); err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "gasto eliminado"})
}

func mapGastoToResponse(g *domain.Gasto) gastoResponse {
	return gastoResponse{
		ID:          g.ID,
		EmpresaID:   g.EmpresaID,
		Monto:       g.Monto,
		Fecha:       g.Fecha,
		TipoPagoID:  g.TipoPagoID,
		Descripcion: g.Descripcion,
	}
}
