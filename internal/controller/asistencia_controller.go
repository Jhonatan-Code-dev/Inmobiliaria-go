package controller

import (
	"strconv"
	"time"

	"rentals-go/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type AsistenciaController struct {
	svc domain.AsistenciaService
}

func NewAsistenciaController(svc domain.AsistenciaService) *AsistenciaController {
	return &AsistenciaController{svc: svc}
}

// --- Operaciones del Empleado ---

// MarcarAsistencia godoc
// @Summary Registrar entrada o salida
// @Description Registra automáticamente si es entrada o salida según los registros del día.
// @Tags Asistencia
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.Asistencia
// @Router /api/user/asistencia/marcar [post]
func (h *AsistenciaController) MarcarAsistencia(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	usuarioID := c.Locals("usuario_id").(int)

	asistencia, err := h.svc.MarcarAsistencia(c.Context(), usuarioID, empresaID)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(asistencia)
}

// MiHistorial godoc
// @Summary Ver mi historial de asistencia
// @Tags Asistencia
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.Asistencia
// @Router /api/user/asistencia/mi-historial [get]
func (h *AsistenciaController) MiHistorial(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	usuarioID := c.Locals("usuario_id").(int)

	historial, err := h.svc.ListarMiHistorial(c.Context(), usuarioID, empresaID)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(historial)
}

// SolicitarPermiso godoc
// @Summary Solicitar un permiso o justificación
// @Tags Asistencia
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body domain.RegistroPermiso true "Datos del permiso"
// @Success 201 {object} domain.Permiso
// @Router /api/user/asistencia/permisos [post]
func (h *AsistenciaController) SolicitarPermiso(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	usuarioID := c.Locals("usuario_id").(int)

	var req domain.RegistroPermiso
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	permiso, err := h.svc.SolicitarPermiso(c.Context(), usuarioID, empresaID, &req)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(permiso)
}

// --- Operaciones del Administrador ---

// ListarRegistros godoc
// @Summary Listar asistencia de todo el personal
// @Tags Asistencia
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {array} domain.Asistencia
// @Router /api/user/asistencia/registros [get]
func (h *AsistenciaController) ListarRegistros(c *fiber.Ctx) error {
	empresaID := c.QueryInt("empresa_id")
	if empresaID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "empresa_id es requerido"})
	}

	filtros := domain.AsistenciaFiltros{
		EmpresaID: empresaID,
		UsuarioID: c.QueryInt("usuario_id"),
		Estado:    c.Query("estado"),
		Pagina:    c.QueryInt("pag", 1),
		Limite:    c.QueryInt("limite", 50),
	}

	if d := c.Query("desde"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		filtros.Desde = &t
	}
	if d := c.Query("hasta"); d != "" {
		t, _ := time.Parse("2006-01-02", d)
		filtros.Hasta = &t
	}

	lista, _, err := h.svc.ListarAsistencia(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(lista)
}

// AsignarHorario godoc
// @Summary Asignar horario a un trabajador
// @Tags Asistencia
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Param request body domain.RegistroHorario true "Datos del horario"
// @Success 200 {object} domain.Horario
// @Router /api/user/asistencia/horarios [post]
func (h *AsistenciaController) AsignarHorario(c *fiber.Ctx) error {
	empresaID := c.QueryInt("empresa_id")
	if empresaID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "empresa_id es requerido"})
	}

	var req domain.RegistroHorario
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	horario, err := h.svc.AsignarHorario(c.Context(), empresaID, &req)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(horario)
}

// DecidirPermiso godoc
// @Summary Aprobar o rechazar un permiso
// @Tags Asistencia
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del permiso"
// @Param empresa_id query int true "ID de la empresa"
// @Param request body domain.DecisionPermiso true "Decisión"
// @Success 200 {object} domain.Permiso
// @Router /api/user/asistencia/permisos/{id}/estado [put]
func (h *AsistenciaController) DecidirPermiso(c *fiber.Ctx) error {
	empresaID := c.QueryInt("empresa_id")
	if empresaID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "empresa_id es requerido"})
	}

	id, _ := strconv.Atoi(c.Params("id"))

	var req domain.DecisionPermiso
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	permiso, err := h.svc.DecidirPermiso(c.Context(), id, empresaID, &req)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(permiso)
}

// EliminarRegistro godoc
// @Summary Eliminar un registro de asistencia
// @Tags Asistencia
// @Security BearerAuth
// @Param id path int true "ID de la asistencia"
// @Param empresa_id query int true "ID de la empresa"
// @Success 204 "No Content"
// @Router /api/user/asistencia/registros/{id} [delete]
func (h *AsistenciaController) EliminarRegistro(c *fiber.Ctx) error {
	empresaID := c.QueryInt("empresa_id")
	if empresaID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "empresa_id es requerido"})
	}

	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.svc.EliminarAsistencia(c.Context(), id, empresaID); err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.SendStatus(204)
}

// ObtenerHorario godoc
// @Summary Obtener horario de un trabajador
// @Tags Asistencia
// @Security BearerAuth
// @Param usuario_id query int true "ID del usuario"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} domain.Horario
// @Router /api/user/asistencia/horarios/detalle [get]
func (h *AsistenciaController) ObtenerHorario(c *fiber.Ctx) error {
	empresaID := c.QueryInt("empresa_id")
	usuarioID := c.QueryInt("usuario_id")
	if empresaID <= 0 || usuarioID <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "empresa_id y usuario_id son requeridos"})
	}

	horario, err := h.svc.ObtenerHorario(c.Context(), usuarioID, empresaID)
	if err != nil {
		return c.Status(404).JSON(errorResponse{Message: "horario no encontrado"})
	}

	return c.JSON(horario)
}
