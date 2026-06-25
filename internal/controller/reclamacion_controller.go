package controller

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"rentals-go/internal/domain"

	"github.com/gofiber/fiber/v2"
)

var (
	rxNombre   = regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s'-]+$`)
	rxTelefono = regexp.MustCompile(`^[+]*[(]{0,1}[0-9]{1,4}[)]{0,1}[-\s./0-9]{6,15}$`)
	rxEmail    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	rxDNI      = regexp.MustCompile(`^[0-9]{8}$`)
	rxAlfanum  = regexp.MustCompile(`^[a-zA-Z0-9-]{5,15}$`)
)

type ReclamacionController struct {
	svc     domain.ReclamacionService
	empRepo domain.EmpresaRepository
}

func NewReclamacionController(svc domain.ReclamacionService, empRepo domain.EmpresaRepository) *ReclamacionController {
	return &ReclamacionController{
		svc:     svc,
		empRepo: empRepo,
	}
}

type registrarReclamacionRequest struct {
	EmpresaID          int     `json:"empresa_id"`
	Nombres            string  `json:"nombres"`
	Apellidos          string  `json:"apellidos"`
	TipoDocumento      string  `json:"tipo_documento"`
	NumeroDocumento    string  `json:"numero_documento"`
	Telefono           string  `json:"telefono"`
	Email              string  `json:"email"`
	Direccion          string  `json:"direccion"`
	MenorEdad          bool    `json:"menor_edad"`
	NombreApoderado    string  `json:"nombre_apoderado"`
	TipoBien           string  `json:"tipo_bien"`
	MontoReclamado     float64 `json:"monto_reclamado"`
	DescripcionBien    string  `json:"descripcion_bien"`
	TipoReclamacion    string  `json:"tipo_reclamacion"`
	DetalleReclamacion string  `json:"detalle_reclamacion"`
	PedidoConsumidor   string  `json:"pedido_consumidor"`
}

type respuestaReclamacionRequest struct {
	Respuesta string `json:"respuesta"`
}

type listadoReclamacionesResponse struct {
	Datos      []*domain.Reclamacion `json:"datos"`
	Paginacion paginadorResponse     `json:"paginacion"`
}

type empresaPublicaResponse struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

func (h *ReclamacionController) ListarPublicasEmpresas(c *fiber.Ctx) error {
	empresas, _, err := h.empRepo.ListarPaginado(c.Context(), 1000, 0, "")
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	out := make([]empresaPublicaResponse, 0)
	for _, emp := range empresas {
		if emp.Estado {
			out = append(out, empresaPublicaResponse{
				ID:     emp.ID,
				Nombre: emp.Nombre,
			})
		}
	}
	return c.JSON(out)
}

func (h *ReclamacionController) RegistrarPublica(c *fiber.Ctx) error {
	var req registrarReclamacionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de petición inválido"})
	}

	if err := validarReclamacion(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: err.Error()})
	}

	rec := &domain.Reclamacion{
		EmpresaID:          req.EmpresaID,
		Nombres:            req.Nombres,
		Apellidos:          req.Apellidos,
		TipoDocumento:      req.TipoDocumento,
		NumeroDocumento:    req.NumeroDocumento,
		Telefono:           req.Telefono,
		Email:              req.Email,
		Direccion:          req.Direccion,
		MenorEdad:          req.MenorEdad,
		NombreApoderado:    req.NombreApoderado,
		TipoBien:           req.TipoBien,
		MontoReclamado:     req.MontoReclamado,
		DescripcionBien:    req.DescripcionBien,
		TipoReclamacion:    req.TipoReclamacion,
		DetalleReclamacion: req.DetalleReclamacion,
		PedidoConsumidor:   req.PedidoConsumidor,
	}

	nuevo, err := h.svc.Registrar(c.Context(), rec)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.Status(201).JSON(nuevo)
}

func validarReclamacion(req *registrarReclamacionRequest) error {
	if req.EmpresaID <= 0 {
		return fmt.Errorf("Debe seleccionar una empresa válida")
	}

	req.Nombres = strings.TrimSpace(req.Nombres)
	req.Apellidos = strings.TrimSpace(req.Apellidos)
	req.NumeroDocumento = strings.TrimSpace(req.NumeroDocumento)
	req.Telefono = strings.TrimSpace(req.Telefono)
	req.Email = strings.TrimSpace(req.Email)
	req.Direccion = strings.TrimSpace(req.Direccion)
	req.NombreApoderado = strings.TrimSpace(req.NombreApoderado)
	req.DescripcionBien = strings.TrimSpace(req.DescripcionBien)
	req.DetalleReclamacion = strings.TrimSpace(req.DetalleReclamacion)
	req.PedidoConsumidor = strings.TrimSpace(req.PedidoConsumidor)

	if req.Nombres == "" {
		return fmt.Errorf("El nombre es obligatorio")
	}
	if len(req.Nombres) < 2 || len(req.Nombres) > 150 || !rxNombre.MatchString(req.Nombres) {
		return fmt.Errorf("El nombre es inválido (solo letras, mínimo 2 y máximo 150 caracteres)")
	}

	if req.Apellidos == "" {
		return fmt.Errorf("El apellido es obligatorio")
	}
	if len(req.Apellidos) < 2 || len(req.Apellidos) > 150 || !rxNombre.MatchString(req.Apellidos) {
		return fmt.Errorf("El apellido es inválido (solo letras, mínimo 2 y máximo 150 caracteres)")
	}

	if req.TipoDocumento != "DNI" && req.TipoDocumento != "CE" && req.TipoDocumento != "PASAPORTE" {
		return fmt.Errorf("Tipo de documento inválido (debe ser DNI, CE o PASAPORTE)")
	}

	if req.NumeroDocumento == "" {
		return fmt.Errorf("El número de documento es obligatorio")
	}
	if req.TipoDocumento == "DNI" && !rxDNI.MatchString(req.NumeroDocumento) {
		return fmt.Errorf("El DNI debe tener exactamente 8 dígitos numéricos")
	}
	if (req.TipoDocumento == "CE" || req.TipoDocumento == "PASAPORTE") && !rxAlfanum.MatchString(req.NumeroDocumento) {
		return fmt.Errorf("El número de documento es inválido (debe ser alfanumérico entre 5 y 15 caracteres)")
	}

	if req.Telefono == "" {
		return fmt.Errorf("El teléfono es obligatorio")
	}
	if !rxTelefono.MatchString(req.Telefono) {
		return fmt.Errorf("El teléfono ingresado es inválido (debe tener entre 7 y 15 dígitos)")
	}

	if req.Email == "" {
		return fmt.Errorf("El correo electrónico es obligatorio")
	}
	if len(req.Email) > 100 || !rxEmail.MatchString(req.Email) {
		return fmt.Errorf("El correo electrónico ingresado no tiene un formato válido")
	}

	if req.Direccion == "" {
		return fmt.Errorf("La dirección de domicilio es obligatoria")
	}
	if len(req.Direccion) < 5 || len(req.Direccion) > 255 {
		return fmt.Errorf("La dirección debe tener entre 5 y 255 caracteres")
	}

	if req.MenorEdad {
		if req.NombreApoderado == "" {
			return fmt.Errorf("El nombre del apoderado es obligatorio para menores de edad")
		}
		if len(req.NombreApoderado) < 5 || len(req.NombreApoderado) > 200 || !rxNombre.MatchString(req.NombreApoderado) {
			return fmt.Errorf("El nombre del apoderado es inválido (solo letras, mínimo 5 caracteres)")
		}
	} else {
		req.NombreApoderado = ""
	}

	if req.TipoBien != "PRODUCTO" && req.TipoBien != "SERVICIO" {
		return fmt.Errorf("Tipo de bien inválido (debe ser PRODUCTO o SERVICIO)")
	}

	if req.MontoReclamado < 0 {
		return fmt.Errorf("El monto reclamado no puede ser negativo")
	}

	if req.DescripcionBien == "" {
		return fmt.Errorf("La descripción del bien es obligatoria")
	}
	if len(req.DescripcionBien) < 5 || len(req.DescripcionBien) > 1000 {
		return fmt.Errorf("La descripción del bien debe tener entre 5 y 1000 caracteres")
	}

	if req.TipoReclamacion != "RECLAMO" && req.TipoReclamacion != "QUEJA" {
		return fmt.Errorf("Tipo de reclamación inválido (debe ser RECLAMO o QUEJA)")
	}

	if req.DetalleReclamacion == "" {
		return fmt.Errorf("El detalle de la reclamación es obligatorio")
	}
	if len(req.DetalleReclamacion) < 5 || len(req.DetalleReclamacion) > 4000 {
		return fmt.Errorf("El detalle de la reclamación debe tener entre 5 y 4000 caracteres")
	}

	if req.PedidoConsumidor == "" {
		return fmt.Errorf("El pedido concreto es obligatorio")
	}
	if len(req.PedidoConsumidor) < 5 || len(req.PedidoConsumidor) > 2000 {
		return fmt.Errorf("El pedido concreto debe tener entre 5 y 2000 caracteres")
	}

	return nil
}

func (h *ReclamacionController) DescargarPDFPublica(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	empresaIDQuery := c.QueryInt("empresa_id")

	pdfBytes, err := h.svc.GenerarPDF(c.Context(), id, empresaIDQuery)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", `attachment; filename="reclamacion_`+strconv.Itoa(id)+`.pdf"`)
	return c.Send(pdfBytes)
}

func (h *ReclamacionController) Listar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	pag := c.QueryInt("pag", 1)
	limite := 10

	list, total, err := h.svc.Listar(c.Context(), empresaID, pag)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	paginas := (total + limite - 1) / limite

	return c.JSON(listadoReclamacionesResponse{
		Datos: list,
		Paginacion: paginadorResponse{
			Total:        total,
			Paginas:      paginas,
			Pagina:       pag,
			PaginaActual: pag,
			PorPagina:    limite,
		},
	})
}

func (h *ReclamacionController) Obtener(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	rec, err := h.svc.Obtener(c.Context(), id, empresaID)
	if err != nil {
		return c.Status(404).JSON(errorResponse{Message: "Reclamación no encontrada"})
	}

	return c.JSON(rec)
}

func (h *ReclamacionController) Responder(c *fiber.Ctx) error {
	empresaID := c.Locals("empresa_id").(int)

	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	var req respuestaReclamacionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(errorResponse{Message: "formato de petición inválido"})
	}

	if req.Respuesta == "" {
		return c.Status(400).JSON(errorResponse{Message: "La respuesta no puede estar vacía"})
	}

	rec, err := h.svc.Responder(c.Context(), id, empresaID, req.Respuesta)
	if err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(rec)
}

func (h *ReclamacionController) Eliminar(c *fiber.Ctx) error {
	empresaID, errResp := obtenerEmpresaIDListado(c)
	if errResp != nil {
		return c.Status(errResp.Code).JSON(errorResponse{Message: errResp.Message})
	}

	id, _ := strconv.Atoi(c.Params("id"))
	if id <= 0 {
		return c.Status(400).JSON(errorResponse{Message: "ID inválido"})
	}

	if err := h.svc.Eliminar(c.Context(), id, empresaID); err != nil {
		return c.Status(500).JSON(errorResponse{Message: err.Error()})
	}

	return c.JSON(fiber.Map{"message": "reclamación eliminada"})
}
