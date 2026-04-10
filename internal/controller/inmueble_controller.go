package controller

import (
	"strconv"

	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/money"

	"github.com/gofiber/fiber/v2"
)

type InmuebleController struct {
	svc domain.InmuebleService
}

func NewInmuebleController(svc domain.InmuebleService) *InmuebleController {
	return &InmuebleController{svc: svc}
}

type inmuebleRequest struct {
	EmpresaID     int     `json:"empresa_id"`
	Nombre        string  `json:"nombre"`
	Tipo          string  `json:"tipo"`
	Descripcion   *string `json:"descripcion"`
	Direccion     string  `json:"direccion"`
	Ciudad        *string `json:"ciudad"`
	Region        *string `json:"region"`
	Pais          *string `json:"pais"`
	CodigoPostal  *string `json:"codigo_postal"`
	TotalPisos    int     `json:"total_pisos"`
	TotalUnidades int     `json:"total_unidades"`
	Estado        string  `json:"estado"`
}

type unidadRequest struct {
	Codigo            string   `json:"codigo"`
	Nombre            *string  `json:"nombre"`
	Tipo              string   `json:"tipo"`
	NumeroPiso        *int     `json:"numero_piso"`
	Dormitorios       int      `json:"dormitorios"`
	Banos             int      `json:"banos"`
	AreaM2            *float64 `json:"area_m2"`
	Capacidad         int      `json:"capacidad"`
	Moneda            string       `json:"moneda"`
	PrecioBase        money.Amount `json:"precio_base" swaggertype:"number" example:"850.00"`
	DepositoRequerido money.Amount `json:"deposito_requerido" swaggertype:"number" example:"500.00"`
	IncluyeAgua       bool     `json:"incluye_agua"`
	IncluyeLuz        bool     `json:"incluye_luz"`
	IncluyeInternet   bool     `json:"incluye_internet"`
	Notas             *string  `json:"notas"`
	Estado            string   `json:"estado"`
}

type inmuebleResponse struct {
	ID            int              `json:"id"`
	EmpresaID     int              `json:"empresa_id"`
	Nombre        string           `json:"nombre"`
	Tipo          string           `json:"tipo"`
	Descripcion   *string          `json:"descripcion"`
	Direccion     string           `json:"direccion"`
	Ciudad        *string          `json:"ciudad"`
	Region        *string          `json:"region"`
	Pais          *string          `json:"pais"`
	CodigoPostal  *string          `json:"codigo_postal"`
	TotalPisos    int              `json:"total_pisos"`
	TotalUnidades int              `json:"total_unidades"`
	Estado        string           `json:"estado"`
	CreadoEn      string           `json:"creado_en"`
	Unidades      []unidadResponse `json:"unidades,omitempty"`
}

type unidadResponse struct {
	ID                 int      `json:"id"`
	PropiedadID        int      `json:"propiedad_id"`
	Codigo             string   `json:"codigo"`
	Nombre             *string  `json:"nombre"`
	Tipo               string   `json:"tipo"`
	NumeroPiso         *int     `json:"numero_piso"`
	Dormitorios        int      `json:"dormitorios"`
	Banos              int      `json:"banos"`
	AreaM2             *float64 `json:"area_m2"`
	Capacidad          int      `json:"capacidad"`
	Moneda             string       `json:"moneda"`
	PrecioBase         money.Amount `json:"precio_base" swaggertype:"number" example:"850.00"`
	DepositoRequerido  money.Amount `json:"deposito_requerido" swaggertype:"number" example:"500.00"`
	IncluyeAgua        bool     `json:"incluye_agua"`
	IncluyeLuz         bool     `json:"incluye_luz"`
	IncluyeInternet    bool     `json:"incluye_internet"`
	Notas              *string  `json:"notas"`
	Estado             string   `json:"estado"`
	CreadoEn           string       `json:"creado_en"`
}

type listadoInmueblesResponse struct {
	Datos      []inmuebleResponse `json:"datos"`
	Paginacion paginadorResponse  `json:"paginacion"`
}

// ListarInmuebles godoc
// @Summary Listar inmuebles con filtros y paginación
// @Description Obtiene una lista paginada de inmuebles (máx 10) filtrados por empresa. Permite buscar por nombre, dirección o ciudad.
// @Tags Inmuebles
// @Security BearerAuth
// @Produce json
// @Param empresa_id query int true "ID de la empresa"
// @Param pag query int false "Número de página" default(1)
// @Param buscar query string false "Búsqueda por nombre, dirección o ciudad"
// @Param estado query string false "Filtrar por estado"
// @Param tipo query string false "Filtrar por tipo de inmueble"
// @Success 200 {object} listadoInmueblesResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/inmuebles [get]
func (h *InmuebleController) Listar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	filtros := domain.InmuebleFiltros{
		EmpresaID: empresaID,
		Busqueda:  c.Query("buscar"),
		Pagina:    c.QueryInt("pag", 1),
		Limite:    10,
		Estado:    c.Query("estado"),
		Tipo:      c.Query("tipo"),
	}

	list, total, err := h.svc.Listar(c.Context(), filtros)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	datos := make([]inmuebleResponse, 0, len(list))
	for _, item := range list {
		datos = append(datos, mapInmuebleResponse(item))
	}

	paginas := (total + filtros.Limite - 1) / filtros.Limite

	return c.JSON(listadoInmueblesResponse{
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

// ObtenerInmueble godoc
// @Summary Obtener inmueble por ID
// @Description Devuelve el detalle del inmueble con sus unidades. Requiere enviar empresa_id por query string.
// @Tags Inmuebles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} inmuebleResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id} [get]
func (h *InmuebleController) Obtener(c *fiber.Ctx) error {
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
		return manejarErrorCliente(err, "inmueble")
	}

	return c.JSON(mapInmuebleResponse(item))
}

// CrearInmueble godoc
// @Summary Registrar un nuevo inmueble
// @Description Crea un nuevo inmueble para la empresa.
// @Tags Inmuebles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body inmuebleRequest true "Datos del inmueble"
// @Success 201 {object} inmuebleResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Router /api/user/inmuebles [post]
func (h *InmuebleController) Crear(c *fiber.Ctx) error {
	empresaID, _ := c.Locals("empresa_id").(int)
	var req inmuebleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	finalEmpresaID, errResp := validarEmpresaIDConSesion(req.EmpresaID, empresaID)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	item := &domain.Inmueble{
		EmpresaID:     finalEmpresaID,
		Nombre:        req.Nombre,
		Tipo:          req.Tipo,
		Descripcion:   req.Descripcion,
		Direccion:     req.Direccion,
		Ciudad:        req.Ciudad,
		Region:        req.Region,
		Pais:          req.Pais,
		CodigoPostal:  req.CodigoPostal,
		TotalPisos:    req.TotalPisos,
		TotalUnidades: req.TotalUnidades,
		Estado:        req.Estado,
	}

	created, err := h.svc.Crear(c.Context(), item)
	if err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}
	return c.Status(201).JSON(mapInmuebleResponse(created))
}

// ActualizarInmueble godoc
// @Summary Actualizar un inmueble
// @Description Actualiza los datos del inmueble existente.
// @Tags Inmuebles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param request body inmuebleRequest true "Datos del inmueble"
// @Success 200 {object} inmuebleResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id} [put]
func (h *InmuebleController) Actualizar(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	var req inmuebleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	finalEmpresaID, errResp := validarEmpresaIDConSesion(req.EmpresaID, empresaID)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	item := &domain.Inmueble{
		ID:            id,
		EmpresaID:     finalEmpresaID,
		Nombre:        req.Nombre,
		Tipo:          req.Tipo,
		Descripcion:   req.Descripcion,
		Direccion:     req.Direccion,
		Ciudad:        req.Ciudad,
		Region:        req.Region,
		Pais:          req.Pais,
		CodigoPostal:  req.CodigoPostal,
		TotalPisos:    req.TotalPisos,
		TotalUnidades: req.TotalUnidades,
		Estado:        req.Estado,
	}

	updated, err := h.svc.Actualizar(c.Context(), item)
	if err != nil {
		return manejarErrorCliente(err, "inmueble")
	}
	return c.JSON(mapInmuebleResponse(updated))
}

// EliminarInmueble godoc
// @Summary Eliminar un inmueble
// @Description Elimina el inmueble. Requiere enviar empresa_id por query string.
// @Tags Inmuebles
// @Security BearerAuth
// @Param id path int true "ID del inmueble"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id} [delete]
func (h *InmuebleController) Eliminar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	id, _ := strconv.Atoi(c.Params("id"))
	if id == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return manejarErrorCliente(err, "inmueble")
	}
	return c.JSON(fiber.Map{"message": "inmueble eliminado"})
}

// ListarUnidadesInmueble godoc
// @Summary Listar unidades de un inmueble
// @Description Devuelve todas las unidades del inmueble. Requiere enviar empresa_id por query string.
// @Tags Inmuebles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {array} unidadResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id}/unidades [get]
func (h *InmuebleController) ListarUnidades(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	propiedadID, _ := strconv.Atoi(c.Params("id"))
	if propiedadID == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	list, err := h.svc.ListarUnidades(c.Context(), propiedadID, empresaID)
	if err != nil {
		return manejarErrorCliente(err, "inmueble")
	}
	out := make([]unidadResponse, 0, len(list))
	for _, item := range list {
		out = append(out, mapUnidadResponse(item))
	}
	return c.JSON(out)
}

// CrearUnidadInmueble godoc
// @Summary Registrar unidad en un inmueble
// @Description Crea una nueva unidad asociada a un inmueble existente.
// @Tags Inmuebles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param request body unidadRequest true "Datos de la unidad"
// @Success 201 {object} unidadResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id}/unidades [post]
func (h *InmuebleController) CrearUnidad(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	propiedadID, _ := strconv.Atoi(c.Params("id"))
	if propiedadID == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	var req unidadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	item := &domain.Unidad{
		Codigo:            req.Codigo,
		Nombre:            req.Nombre,
		Tipo:              req.Tipo,
		NumeroPiso:        req.NumeroPiso,
		Dormitorios:       req.Dormitorios,
		Banos:             req.Banos,
		AreaM2:            req.AreaM2,
		Capacidad:         req.Capacidad,
		Moneda:            req.Moneda,
		PrecioBase:        req.PrecioBase.Float64(),
		PrecioBaseCents:   req.PrecioBase.Cents(),
		DepositoRequerido: req.DepositoRequerido.Float64(),
		DepositoReqCents:  req.DepositoRequerido.Cents(),
		IncluyeAgua:       req.IncluyeAgua,
		IncluyeLuz:        req.IncluyeLuz,
		IncluyeInternet:   req.IncluyeInternet,
		Notas:             req.Notas,
		Estado:            req.Estado,
	}
	created, err := h.svc.CrearUnidad(c.Context(), propiedadID, empresaID, item)
	if err != nil {
		return manejarErrorCliente(err, "unidad")
	}
	return c.Status(201).JSON(mapUnidadResponse(created))
}

// ObtenerUnidadInmueble godoc
// @Summary Obtener una unidad
// @Description Devuelve el detalle de una unidad del inmueble. Requiere enviar empresa_id por query string.
// @Tags Inmuebles
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param unidadId path int true "ID de la unidad"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} unidadResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id}/unidades/{unidadId} [get]
func (h *InmuebleController) ObtenerUnidad(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	propiedadID, _ := strconv.Atoi(c.Params("id"))
	unidadID, _ := strconv.Atoi(c.Params("unidadId"))
	if propiedadID == 0 || unidadID == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	item, err := h.svc.ObtenerUnidad(c.Context(), propiedadID, unidadID, empresaID)
	if err != nil {
		return manejarErrorCliente(err, "unidad")
	}
	return c.JSON(mapUnidadResponse(item))
}

// ActualizarUnidadInmueble godoc
// @Summary Actualizar una unidad
// @Description Actualiza los datos de una unidad existente.
// @Tags Inmuebles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID del inmueble"
// @Param unidadId path int true "ID de la unidad"
// @Param request body unidadRequest true "Datos de la unidad"
// @Success 200 {object} unidadResponse
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id}/unidades/{unidadId} [put]
func (h *InmuebleController) ActualizarUnidad(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)
	propiedadID, _ := strconv.Atoi(c.Params("id"))
	unidadID, _ := strconv.Atoi(c.Params("unidadId"))
	if propiedadID == 0 || unidadID == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	var req unidadRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato inválido"})
	}
	item := &domain.Unidad{
		ID:                unidadID,
		Codigo:            req.Codigo,
		Nombre:            req.Nombre,
		Tipo:              req.Tipo,
		NumeroPiso:        req.NumeroPiso,
		Dormitorios:       req.Dormitorios,
		Banos:             req.Banos,
		AreaM2:            req.AreaM2,
		Capacidad:         req.Capacidad,
		Moneda:            req.Moneda,
		PrecioBase:        req.PrecioBase.Float64(),
		PrecioBaseCents:   req.PrecioBase.Cents(),
		DepositoRequerido: req.DepositoRequerido.Float64(),
		DepositoReqCents:  req.DepositoRequerido.Cents(),
		IncluyeAgua:       req.IncluyeAgua,
		IncluyeLuz:        req.IncluyeLuz,
		IncluyeInternet:   req.IncluyeInternet,
		Notas:             req.Notas,
		Estado:            req.Estado,
	}
	updated, err := h.svc.ActualizarUnidad(c.Context(), propiedadID, empresaID, item)
	if err != nil {
		return manejarErrorCliente(err, "unidad")
	}
	return c.JSON(mapUnidadResponse(updated))
}

// EliminarUnidadInmueble godoc
// @Summary Eliminar una unidad
// @Description Elimina una unidad del inmueble. Requiere enviar empresa_id por query string.
// @Tags Inmuebles
// @Security BearerAuth
// @Param id path int true "ID del inmueble"
// @Param unidadId path int true "ID de la unidad"
// @Param empresa_id query int true "ID de la empresa"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/user/inmuebles/{id}/unidades/{unidadId} [delete]
func (h *InmuebleController) EliminarUnidad(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}
	propiedadID, _ := strconv.Atoi(c.Params("id"))
	unidadID, _ := strconv.Atoi(c.Params("unidadId"))
	if propiedadID == 0 || unidadID == 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}
	if err := h.svc.EliminarUnidad(c.Context(), propiedadID, unidadID, empresaID); err != nil {
		return manejarErrorCliente(err, "unidad")
	}
	return c.JSON(fiber.Map{"message": "unidad eliminada"})
}

func mapInmuebleResponse(item *domain.Inmueble) inmuebleResponse {
	resp := inmuebleResponse{
		ID:            item.ID,
		EmpresaID:     item.EmpresaID,
		Nombre:        item.Nombre,
		Tipo:          item.Tipo,
		Descripcion:   item.Descripcion,
		Direccion:     item.Direccion,
		Ciudad:        item.Ciudad,
		Region:        item.Region,
		Pais:          item.Pais,
		CodigoPostal:  item.CodigoPostal,
		TotalPisos:    item.TotalPisos,
		TotalUnidades: item.TotalUnidades,
		Estado:        item.Estado,
		CreadoEn:      item.CreadoEn.Format("2006-01-02T15:04:05Z07:00"),
	}
	if len(item.Unidades) > 0 {
		resp.Unidades = make([]unidadResponse, 0, len(item.Unidades))
		for _, unidad := range item.Unidades {
			resp.Unidades = append(resp.Unidades, mapUnidadResponse(unidad))
		}
	}
	return resp
}

func mapUnidadResponse(item *domain.Unidad) unidadResponse {
	return unidadResponse{
		ID:                item.ID,
		PropiedadID:       item.PropiedadID,
		Codigo:            item.Codigo,
		Nombre:            item.Nombre,
		Tipo:              item.Tipo,
		NumeroPiso:        item.NumeroPiso,
		Dormitorios:       item.Dormitorios,
		Banos:             item.Banos,
		AreaM2:            item.AreaM2,
		Capacidad:         item.Capacidad,
		Moneda:            item.Moneda,
		PrecioBase:        money.NewAmountFromCents(item.PrecioBaseCents),
		DepositoRequerido: money.NewAmountFromCents(item.DepositoReqCents),
		IncluyeAgua:       item.IncluyeAgua,
		IncluyeLuz:        item.IncluyeLuz,
		IncluyeInternet:   item.IncluyeInternet,
		Notas:             item.Notas,
		Estado:            item.Estado,
		CreadoEn:          item.CreadoEn.Format("2006-01-02T15:04:05Z07:00"),
	}
}
