package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/money"

	"github.com/gofiber/fiber/v2"
)

type AlquilerController struct {
	svc      domain.AlquilerService
	pagoSvc  domain.PagoAlquilerService
}

func NewAlquilerController(svc domain.AlquilerService, pagoSvc domain.PagoAlquilerService) *AlquilerController {
	return &AlquilerController{svc: svc, pagoSvc: pagoSvc}
}

type alquilerRequest struct {
	ClienteID           int          `json:"cliente_id"`
	UnidadID            int          `json:"unidad_id"`
	FechaInicio         string       `json:"fecha_inicio"`
	FechaFin            *string      `json:"fecha_fin"`
	VencimientoDiaPago  int          `json:"vencimiento_dia_pago"`
	MontoRenta          money.Amount `json:"monto_renta" swaggertype:"number" example:"1200.50"`
	DepositoGarantia    money.Amount `json:"deposito_garantia" swaggertype:"number" example:"1200.50"`
	Moneda              string       `json:"moneda"`
	Observaciones       *string      `json:"observaciones"`
}

type alquilerResponse struct {
	ID          int          `json:"id"`
	ClienteID   int          `json:"cliente_id"`
	Cliente     string       `json:"cliente"`
	UnidadID    int          `json:"unidad_id"`
	Unidad      string       `json:"unidad"`
	Monto       money.Amount `json:"monto"`
	FechaInicio string       `json:"fecha_inicio"`
	FechaFin    *string      `json:"fecha_fin,omitempty"`
	Estado      string       `json:"estado"`
	Moneda      string       `json:"moneda"`
}

type listadoAlquileresResponse struct {
	Datos      []alquilerResponse `json:"datos"`
	Paginacion paginadorResponse  `json:"paginacion"`
}

type registrarPagoRequest struct {
	AlquilerID         int          `json:"alquiler_id"`
	MontoPagado        money.Amount `json:"monto_pagado" swaggertype:"number" example:"1200.50"`
	FechaPago          string       `json:"fecha_pago"`
	MetodoPago         string       `json:"metodo_pago"`
	Nota               *string      `json:"nota"`
	MesCorrespondiente int          `json:"mes_correspondiente"`
}

type pagoAlquilerResponse struct {
	ID                 int          `json:"id"`
	AlquilerID         int          `json:"alquiler_id"`
	ClienteID          *int         `json:"cliente_id,omitempty"`
	Cliente            string       `json:"cliente,omitempty"`
	Unidad             string       `json:"unidad,omitempty"`
	NumeroRecibo       string       `json:"numero_recibo"`
	FechaPago          string       `json:"fecha_pago"`
	Moneda             string       `json:"moneda"`
	MontoPagado        money.Amount `json:"monto_pagado"`
	MetodoPago         string       `json:"metodo_pago"`
	Nota               *string      `json:"nota"`
	MesCorrespondiente int          `json:"mes_correspondiente"`
}

type pendientePagoResponse struct {
	AlquilerID      int          `json:"alquiler_id"`
	Cliente         string       `json:"cliente"`
	Unidad          string       `json:"unidad"`
	Monto           money.Amount `json:"monto"`
	FechaVencimiento string      `json:"fecha_vencimiento"`
	Estado          string       `json:"estado"`
}

// ListarAlquileres godoc
// @Summary Listar alquileres
// @Description Lista contratos de alquiler con filtros y paginación.
// @Tags Alquileres
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Param buscar query string false "Buscar por cliente o código de unidad"
// @Param estado query string false "Estado del contrato"
// @Param unidad_id query int false "ID de la unidad"
// @Param pag query int false "Página"
// @Param por_pagina query int false "Tamaño de página"
// @Success 200 {object} listadoAlquileresResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/alquileres [get]
func (h *AlquilerController) Listar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	porPagina := c.QueryInt("por_pagina", 10)
	if porPagina <= 0 {
		porPagina = 10
	}
	filtros := domain.AlquilerFiltros{
		EmpresaID: empresaID,
		Busqueda:  c.Query("buscar"),
		Estado:    c.Query("estado"),
		UnidadID:  c.QueryInt("unidad_id"),
		Pagina:    c.QueryInt("pag", 1),
		Limite:    porPagina,
	}

	list, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	datos := make([]alquilerResponse, 0, len(list))
	for _, item := range list {
		datos = append(datos, mapAlquilerResponse(item))
	}
	paginas := 0
	if filtros.Limite > 0 {
		paginas = (total + filtros.Limite - 1) / filtros.Limite
	}

	return c.JSON(listadoAlquileresResponse{
		Datos: datos,
		Paginacion: paginadorResponse{
			Total:        total,
			Paginas:      paginas,
			Pagina:       filtros.Pagina,
			PaginaActual: filtros.Pagina,
			PorPagina:    filtros.Limite,
		},
	})
}

// ListarActivosSelector godoc
// @Summary Listar contratos activos para selectores
// @Description Retorna una lista simplificada de contratos activos.
// @Tags Alquileres
// @Router /api/user/alquileres/activos/selector [get]
func (h *AlquilerController) ListarActivosSelector(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	list, err := h.svc.ListarActivosSelector(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	
	datos := make([]alquilerResponse, 0, len(list))
	for _, item := range list {
		datos = append(datos, mapAlquilerResponse(item))
	}
	return c.JSON(datos)
}

// ObtenerAlquiler godoc
// @Summary Obtener alquiler por ID
// @Description Devuelve el detalle de un contrato de alquiler.
// @Tags Alquileres
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del alquiler"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} alquilerResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/alquileres/{id} [get]
func (h *AlquilerController) Obtener(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	item, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return manejarErrorCliente(err, "alquiler")
	}
	return c.JSON(mapAlquilerResponse(item))
}

// CrearAlquiler godoc
// @Summary Registrar un alquiler
// @Description Crea un nuevo contrato y ocupa la unidad.
// @Tags Alquileres
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body alquilerRequest true "Datos del contrato"
// @Success 201 {object} alquilerResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/alquileres [post]
func (h *AlquilerController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req alquilerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicio)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "fecha_inicio debe tener formato YYYY-MM-DD"})
	}

	var fechaFin *time.Time
	if req.FechaFin != nil && *req.FechaFin != "" {
		t, err := time.Parse("2006-01-02", *req.FechaFin)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "fecha_fin debe tener formato YYYY-MM-DD"})
		}
		fechaFin = &t
	}

	item := &domain.Alquiler{
		EmpresaID:        empresaID,
		ClienteID:        req.ClienteID,
		UnidadID:         req.UnidadID,
		FechaInicio:      fechaInicio,
		FechaFin:         fechaFin,
		DiaVencimiento:   req.VencimientoDiaPago,
		Moneda:           req.Moneda,
		MontoRenta:       req.MontoRenta.Float64(),
		MontoRentaCents:  req.MontoRenta.Cents(),
		MontoDeposito:    req.DepositoGarantia.Float64(),
		MontoDepositoCts: req.DepositoGarantia.Cents(),
		Observaciones:    req.Observaciones,
	}

	created, err := h.svc.Crear(c.Context(), item)
	if err != nil {
		return manejarErrorCliente(err, "alquiler")
	}
	return c.Status(201).JSON(mapAlquilerResponse(created))
}

// ActualizarAlquiler godoc
// @Summary Editar contrato de alquiler
// @Tags Alquileres
// @Param id path int true "ID del alquiler"
// @Param request body alquilerRequest true "Datos a actualizar"
// @Router /api/user/alquileres/{id} [put]
func (h *AlquilerController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	var req alquilerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	
	fechaInicio, _ := time.Parse("2006-01-02", req.FechaInicio)
	var fechaFin *time.Time
	if req.FechaFin != nil && *req.FechaFin != "" {
		t, _ := time.Parse("2006-01-02", *req.FechaFin)
		fechaFin = &t
	}

	item := &domain.Alquiler{
		FechaInicio:      fechaInicio,
		FechaFin:         fechaFin,
		DiaVencimiento:   req.VencimientoDiaPago,
		Moneda:           req.Moneda,
		MontoRentaCents:  req.MontoRenta.Cents(),
		MontoDepositoCts: req.DepositoGarantia.Cents(),
		Observaciones:    req.Observaciones,
	}

	updated, err := h.svc.Actualizar(c.Context(), id, empresaID, item)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(mapAlquilerResponse(updated))
}

// EliminarAlquiler godoc
// @Summary Eliminar contrato de alquiler
// @Tags Alquileres
// @Param id path int true "ID del alquiler"
// @Router /api/user/alquileres/{id} [delete]
func (h *AlquilerController) Eliminar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		if strings.Contains(err.Error(), "foreign key constraint fails") {
			return c.Status(400).JSON(errorResponse{Message: "No se puede eliminar el alquiler porque tiene pagos o cargos en su historial."})
		}
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(fiber.Map{"message": "contrato eliminado correctamente"})
}

// TerminarContrato godoc
// @Summary Finalizar contrato formalmente
// @Tags Alquileres
// @Param id path int true "ID del alquiler"
// @Router /api/user/alquileres/{id}/terminar [post]
func (h *AlquilerController) TerminarContrato(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	if err := h.svc.Terminar(c.Context(), id, empresaID); err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(fiber.Map{"message": "contrato finalizado correctamente"})
}

// RegistrarPago godoc
// @Summary Registrar pago de alquiler
// @Description Registra un cobro para un alquiler.
// @Tags Pagos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body registrarPagoRequest true "Datos del pago"
// @Success 201 {object} pagoAlquilerResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/pagos [post]
func (h *AlquilerController) RegistrarPago(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req registrarPagoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	fechaPago, err := time.Parse("2006-01-02", req.FechaPago)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: "fecha_pago debe tener formato YYYY-MM-DD"})
	}
	item := &domain.RegistroPagoAlquiler{
		EmpresaID:          empresaID,
		ContratoID:         req.AlquilerID,
		MontoPagado:        req.MontoPagado.Float64(),
		MontoPagadoCents:   req.MontoPagado.Cents(),
		FechaPago:          fechaPago,
		MetodoPago:         req.MetodoPago,
		Nota:               req.Nota,
		MesCorrespondiente: req.MesCorrespondiente,
	}
	created, err := h.pagoSvc.Registrar(c.Context(), item)
	if err != nil {
		return manejarErrorCliente(err, "pago")
	}
	return c.Status(201).JSON(mapPagoAlquilerResponse(created))
}

// PendientesPago godoc
// @Summary Listar pagos pendientes del mes actual
// @Description Retorna contratos sin pago completo registrado para el mes actual.
// @Tags Pagos
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {array} pendientePagoResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/pagos/pendientes [get]
func (h *AlquilerController) PendientesPago(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	list, err := h.pagoSvc.ListarPendientesMesActual(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	out := make([]pendientePagoResponse, 0, len(list))
	for _, item := range list {
		out = append(out, pendientePagoResponse{
			AlquilerID:       item.AlquilerID,
			Cliente:          item.Cliente,
			Unidad:           item.Unidad,
			Monto:            money.NewAmountFromCents(item.MontoCents),
			FechaVencimiento: item.FechaVencimiento.Format("2006-01-02"),
			Estado:           item.Estado,
		})
	}
	return c.JSON(out)
}

// ListarPagos godoc
// @Summary Historial de pagos
// @Description Lista pagos realizados con filtros y paginación.
// @Tags Pagos
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Param buscar query string false "Búsqueda por cliente o unidad"
// @Param pag query int false "Página"
// @Param por_pagina query int false "Tamaño de página"
// @Success 200 {object} paginatedResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/user/pagos [get]
func (h *AlquilerController) ListarPagos(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	porPagina := c.QueryInt("por_pagina", 10)
	if porPagina <= 0 {
		porPagina = 10
	}
	filtros := domain.PagoFiltros{
		EmpresaID: empresaID,
		Busqueda:  c.Query("buscar"),
		Pagina:    c.QueryInt("pag", 1),
		Limite:    porPagina,
	}

	list, total, err := h.pagoSvc.ListarHistorial(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	out := make([]pagoAlquilerResponse, 0, len(list))
	for _, item := range list {
		out = append(out, mapPagoAlquilerResponse(item))
	}

	return c.JSON(paginatedResponse{
		Datos: out,
		Paginacion: paginadorResponse{
			Total:        total,
			PaginaActual: filtros.Pagina,
			PorPagina:    filtros.Limite,
			Paginas:      (total + filtros.Limite - 1) / filtros.Limite,
		},
	})
}

// ObtenerPago godoc
// @Summary Detalle de un pago
// @Tags Pagos
// @Router /api/user/pagos/{id} [get]
func (h *AlquilerController) ObtenerPago(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	item, err := h.pagoSvc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(404).JSON(errorResponse{Message: "pago no encontrado"})
	}

	return c.JSON(mapPagoAlquilerResponse(item))
}

// ActualizarPago godoc
// @Summary Editar notas o método de pago
// @Tags Pagos
// @Param id path int true "ID del pago"
// @Param request body map[string]interface{} true "Notas o metodo_pago"
// @Router /api/user/pagos/{id} [put]
func (h *AlquilerController) ActualizarPago(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	
	type updateReq struct {
		MetodoPago string  `json:"metodo_pago"`
		Nota       *string `json:"nota"`
	}
	var req updateReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	updated, err := h.pagoSvc.Actualizar(c.Context(), id, empresaID, req.Nota, req.MetodoPago)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(mapPagoAlquilerResponse(updated))
}

// AnularPago godoc
// @Summary Anular un pago
// @Tags Pagos
// @Router /api/user/pagos/{id} [delete]
func (h *AlquilerController) AnularPago(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")

	if err := h.pagoSvc.Anular(c.Context(), id, empresaID); err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "pago anulado correctamente"})
}

// --- Plantillas ---

// ListarPlantillas godoc
// @Summary Listar plantillas de contrato
// @Tags Alquileres
// @Router /api/user/alquileres/plantillas [get]
func (h *AlquilerController) ListarPlantillas(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	list, err := h.svc.ListarPlantillas(c.Context(), empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(list)
}

// GuardarPlantilla godoc
// @Summary Crear o actualizar plantilla
// @Tags Alquileres
// @Router /api/user/alquileres/plantillas [post]
func (h *AlquilerController) GuardarPlantilla(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.PlantillaContrato
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	req.EmpresaID = empresaID
	res, err := h.svc.GuardarPlantilla(c.Context(), &req)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(res)
}

// EliminarPlantilla godoc
// @Summary Eliminar plantilla
// @Tags Alquileres
// @Router /api/user/alquileres/plantillas/{id} [delete]
func (h *AlquilerController) EliminarPlantilla(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	if err := h.svc.EliminarPlantilla(c.Context(), id, empresaID); err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}
	return c.JSON(fiber.Map{"message": "plantilla eliminada"})
}

// --- Generación ---

// GenerarDocumento godoc
// @Summary Generar documento de contrato (Texto/Markdown)
// @Description Genera el contenido del contrato reemplazando las variables dinámicas.
// @Tags Alquileres
// @Param id path int true "ID del alquiler"
// @Param plantilla_id query int false "ID de la plantilla opcional"
// @Router /api/user/alquileres/{id}/generar-documento [get]
func (h *AlquilerController) GenerarDocumento(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	plantillaID := c.QueryInt("plantilla_id", 0)

	texto, err := h.svc.GenerarContrato(c.Context(), id, empresaID, plantillaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{
		"alquiler_id": id,
		"contenido":   texto,
	})
}

// DescargarWord godoc
// @Summary Descargar contrato en formato Word
// @Description Genera y descarga un archivo .doc compatible con Microsoft Word.
// @Tags Alquileres
// @Param id path int true "ID del alquiler"
// @Param plantilla_id query int false "ID de la plantilla opcional"
// @Produce application/vnd.ms-word
// @Router /api/user/alquileres/{id}/descargar-word [get]
func (h *AlquilerController) DescargarWord(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := c.ParamsInt("id")
	plantillaID := c.QueryInt("plantilla_id", 0)

	bytes, err := h.svc.GenerarContratoWord(c.Context(), id, empresaID, plantillaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	c.Set("Content-Type", "application/msword")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"CONTRATO_ARRENDAMIENTO_%d.doc\"", id))

	return c.Send(bytes)
}

// GenerarBorrador godoc
// @Summary Generar contrato rápido en Word
// @Description Genera un documento Word (.doc) al vuelo con datos manuales.
// @Tags Alquileres
// @Accept json
// @Produce application/msword
// @Param request body domain.GenerarBorradorRequest true "Datos para generar el borrador"
// @Router /api/user/alquileres/generar-borrador [post]
func (h *AlquilerController) GenerarBorrador(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req domain.GenerarBorradorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	bytes, err := h.svc.GenerarContratoBorrador(c.Context(), empresaID, req)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	c.Set("Content-Type", "application/msword")
	c.Set("Content-Disposition", "attachment; filename=\"BORRADOR_CONTRATO_ARRENDAMIENTO.doc\"")

	return c.Send(bytes)
}


func mapAlquilerResponse(item *domain.Alquiler) alquilerResponse {
	resp := alquilerResponse{
		ID:          item.ID,
		ClienteID:   item.ClienteID,
		Cliente:     item.ClienteNombre,
		UnidadID:    item.UnidadID,
		Unidad:      item.UnidadCodigo,
		Monto:       money.NewAmountFromCents(item.MontoRentaCents),
		FechaInicio: item.FechaInicio.Format("2006-01-02"),
		Estado:      item.Estado,
		Moneda:      item.Moneda,
	}
	if item.FechaFin != nil {
		s := item.FechaFin.Format("2006-01-02")
		resp.FechaFin = &s
	}
	return resp
}

func mapPagoAlquilerResponse(item *domain.PagoAlquiler) pagoAlquilerResponse {
	return pagoAlquilerResponse{
		ID:                 item.ID,
		AlquilerID:         item.ContratoID,
		ClienteID:          item.ClienteID,
		Cliente:            item.ClienteNombre,
		Unidad:             item.UnidadCodigo,
		NumeroRecibo:       item.NumeroRecibo,
		FechaPago:          item.FechaPago.Format("2006-01-02"),
		Moneda:             item.Moneda,
		MontoPagado:        money.NewAmountFromCents(item.MontoPagadoCents),
		MetodoPago:         item.MetodoPago,
		Nota:               item.Nota,
		MesCorrespondiente: item.MesCorrespondiente,
	}
}
