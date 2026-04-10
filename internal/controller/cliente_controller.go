package controller

import (
	"errors"
	"strings"
	"strconv"
	"time"

	"rentals-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type ClienteController struct {
	svc domain.ClienteService
}

func NewClienteController(svc domain.ClienteService) *ClienteController {
	return &ClienteController{svc: svc}
}

type clienteRequest struct {
	EmpresaID            int     `json:"empresa_id"`
	TipoIdentificacionID int     `json:"tipo_identificacion_id"`
	DocumentoNumero      string  `json:"documento_numero"`
	Nombres              string  `json:"nombres"`
	Apellidos            *string `json:"apellidos"`
	Correo               *string `json:"correo"`
	FechaNacimiento      *string `json:"fecha_nacimiento"` // YYYY-MM-DD
	Nacionalidad         *string `json:"nacionalidad"`
	Direccion            *string `json:"direccion"`
	ContactoEmergencia   *string `json:"contacto_emergencia"`
	TelefonoEmergencia   *string `json:"telefono_emergencia"`
	Notas                *string `json:"notas"`
	Estado               string  `json:"estado"`
}

type clienteResponse struct {
	ID                   int        `json:"id"`
	EmpresaID            int        `json:"empresa_id"`
	TipoIdentificacionID int        `json:"tipo_identificacion_id"`
	DocumentoNumero      string     `json:"documento_numero"`
	Nombres              string     `json:"nombres"`
	Apellidos            *string    `json:"apellidos"`
	Correo               *string    `json:"correo"`
	FechaNacimiento      *time.Time `json:"fecha_nacimiento"`
	Nacionalidad         *string    `json:"nacionalidad"`
	Direccion            *string    `json:"direccion"`
	ContactoEmergencia   *string    `json:"contacto_emergencia"`
	TelefonoEmergencia   *string    `json:"telefono_emergencia"`
	Notas                *string    `json:"notas"`
	Estado               string     `json:"estado"`
	CreadoEn             time.Time  `json:"creado_en"`
}

type listadoClientesResponse struct {
	Datos      []clienteResponse `json:"datos"`
	Paginacion paginadorResponse `json:"paginacion"`
}

// ListarTiposIdentificacion godoc
// @Summary Catálogo de tipos de identificación
// @Description Lista los tipos de identificación activos que el frontend debe usar para enviar tipo_identificacion_id al crear o actualizar clientes.
// @Tags Clientes
// @Security BearerAuth
// @Produce json
// @Success 200 {array} tipoIdentificacionResponse
// @Failure 500 {object} errorResponse
// @Router /api/user/clientes/tipos-identificacion [get]
func (h *ClienteController) ListarTiposIdentificacion(c *fiber.Ctx) error {
	list, err := h.svc.ListarTiposIdentificacion(c.Context())
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	out := make([]tipoIdentificacionResponse, 0, len(list))
	for _, item := range list {
		out = append(out, tipoIdentificacionResponse{
			ID:     item.ID,
			Codigo: item.Codigo,
			Nombre: item.Nombre,
			Pais:   item.Pais,
			Activo: item.Activo,
		})
	}

	return c.JSON(out)
}

// Obtener godoc
// @Summary Obtener cliente por ID
// @Description Devuelve el detalle de un cliente. Requiere enviar empresa_id por query string.
// @Tags Clientes
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del cliente"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} clienteResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/clientes/{id} [get]
func (h *ClienteController) Obtener(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	cliente, err := h.svc.ObtenerCliente(c.Context(), id, empresaID)
	if err != nil {
		return manejarErrorCliente(err, "cliente")
	}

	return c.JSON(mapClienteToResponse(cliente))
}

// Listar godoc
// @Summary Listar clientes con filtros y paginación
// @Description Obtiene una lista paginada de clientes (máx 10). Filtros opcionales: empresa_id, buscar, pag.
// @Tags Clientes
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Param buscar query string false "Búsqueda por nombre, apellido o documento"
// @Param pag query int false "Número de página" default(1)
// @Success 200 {object} listadoClientesResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/user/clientes [get]
func (h *ClienteController) Listar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	filtros := domain.ClienteFiltros{
		EmpresaID: empresaID,
		Busqueda:  c.Query("buscar"),
		Pagina:    c.QueryInt("pag", 1),
		Limite:    10,
	}

	list, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	datos := make([]clienteResponse, 0, len(list))
	for _, cl := range list {
		datos = append(datos, mapClienteToResponse(cl))
	}

	paginas := (total + filtros.Limite - 1) / filtros.Limite

	return c.JSON(listadoClientesResponse{
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

// Crear godoc
// @Summary Registrar un nuevo cliente
// @Description Crea un nuevo cliente para la empresa.
// @Tags Clientes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body clienteRequest true "Datos del cliente"
// @Success 201 {object} clienteResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/clientes [post]
func (h *ClienteController) Crear(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	var req clienteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	finalEmpresaID, errResp := validarEmpresaIDConSesion(req.EmpresaID, empresaID)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	cliente := &domain.Cliente{
		EmpresaID:            finalEmpresaID,
		TipoIdentificacionID: req.TipoIdentificacionID,
		DocumentoNumero:      req.DocumentoNumero,
		Nombres:              req.Nombres,
		Apellidos:            req.Apellidos,
		Correo:               req.Correo,
		Nacionalidad:         req.Nacionalidad,
		Direccion:            req.Direccion,
		ContactoEmergencia:   req.ContactoEmergencia,
		TelefonoEmergencia:   req.TelefonoEmergencia,
		Notas:                req.Notas,
		Estado:               req.Estado,
	}

	if req.FechaNacimiento != nil && *req.FechaNacimiento != "" {
		t, err := time.Parse("2006-01-02", *req.FechaNacimiento)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "fecha_nacimiento debe tener formato YYYY-MM-DD"})
		}
		cliente.FechaNacimiento = &t
	}

	if cliente.Estado == "" {
		cliente.Estado = "activo"
	}

	nuevo, err := h.svc.RegistrarCliente(c.Context(), cliente)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(mapClienteToResponse(nuevo))
}

// Actualizar godoc
// @Summary Actualizar un cliente
// @Description Actualiza los datos de un cliente existente.
// @Tags Clientes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del cliente"
// @Param request body clienteRequest true "Datos a actualizar"
// @Success 200 {object} clienteResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/clientes/{id} [put]
func (h *ClienteController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	var req clienteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}

	finalEmpresaID, errResp := validarEmpresaIDConSesion(req.EmpresaID, empresaID)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	cliente := &domain.Cliente{
		ID:                   id,
		EmpresaID:            finalEmpresaID,
		TipoIdentificacionID: req.TipoIdentificacionID,
		DocumentoNumero:      req.DocumentoNumero,
		Nombres:              req.Nombres,
		Apellidos:            req.Apellidos,
		Correo:               req.Correo,
		Nacionalidad:         req.Nacionalidad,
		Direccion:            req.Direccion,
		ContactoEmergencia:   req.ContactoEmergencia,
		TelefonoEmergencia:   req.TelefonoEmergencia,
		Notas:                req.Notas,
		Estado:               req.Estado,
	}

	if req.FechaNacimiento != nil && *req.FechaNacimiento != "" {
		t, err := time.Parse("2006-01-02", *req.FechaNacimiento)
		if err != nil {
			return c.Status(400).JSON(errorResponse{Message: "fecha_nacimiento debe tener formato YYYY-MM-DD"})
		}
		cliente.FechaNacimiento = &t
	}

	if cliente.Estado == "" {
		cliente.Estado = "activo"
	}

	actualizado, err := h.svc.ActualizarCliente(c.Context(), cliente)
	if err != nil {
		return manejarErrorCliente(err, "cliente")
	}

	return c.JSON(mapClienteToResponse(actualizado))
}

// Eliminar godoc
// @Summary Eliminar un cliente
// @Description Elimina el registro del cliente.
// @Tags Clientes
// @Security BearerAuth
// @Param id path int true "ID del cliente"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/clientes/{id} [delete]
func (h *ClienteController) Eliminar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	if err := h.svc.EliminarCliente(c.Context(), id, empresaID); err != nil {
		return manejarErrorCliente(err, "cliente")
	}

	return c.JSON(fiber.Map{"message": "cliente eliminado"})
}

func manejarErrorCliente(err error, entidad string) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return fiber.NewError(fiber.StatusNotFound, entidad+" no encontrado")
	case errors.Is(err, domain.ErrForbidden):
		return fiber.NewError(fiber.StatusForbidden, limpiarPrefijoError(err.Error(), "forbidden: "))
	default:
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
}

func limpiarPrefijoError(message, prefix string) string {
	return strings.TrimPrefix(message, prefix)
}

func mapClienteToResponse(c *domain.Cliente) clienteResponse {
	return clienteResponse{
		ID:                   c.ID,
		EmpresaID:            c.EmpresaID,
		TipoIdentificacionID: c.TipoIdentificacionID,
		DocumentoNumero:      c.DocumentoNumero,
		Nombres:              c.Nombres,
		Apellidos:            c.Apellidos,
		Correo:               c.Correo,
		FechaNacimiento:      c.FechaNacimiento,
		Nacionalidad:         c.Nacionalidad,
		Direccion:            c.Direccion,
		ContactoEmergencia:   c.ContactoEmergencia,
		TelefonoEmergencia:   c.TelefonoEmergencia,
		Notas:                c.Notas,
		Estado:               c.Estado,
		CreadoEn:             c.CreadoEn,
	}
}
