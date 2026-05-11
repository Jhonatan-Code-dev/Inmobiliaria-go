package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rentals-go/internal/domain"
)

type AlquilerService struct {
	repo    domain.AlquilerRepository
	cliente domain.ClienteRepository
}

func NewAlquilerService(repo domain.AlquilerRepository, cliente domain.ClienteRepository) *AlquilerService {
	return &AlquilerService{repo: repo, cliente: cliente}
}

func (s *AlquilerService) Listar(ctx context.Context, filtros domain.AlquilerFiltros) ([]*domain.Alquiler, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.ListarPaginado(ctx, filtros)
}

func (s *AlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.Alquiler, error) {
	item, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}
	if item.EmpresaID != empresaID {
		return nil, fmt.Errorf("%w: alquiler no pertenece a la empresa", domain.ErrForbidden)
	}
	return item, nil
}

func (s *AlquilerService) Crear(ctx context.Context, alquiler *domain.Alquiler) (*domain.Alquiler, error) {
	if alquiler.Tipo == "" {
		alquiler.Tipo = "alquiler"
	}
	if alquiler.Moneda == "" {
		alquiler.Moneda = "PEN"
	}
	if alquiler.Estado == "" {
		alquiler.Estado = "activo"
	}
	alquiler.ActivoParaCobro = true
	alquiler.Codigo = fmt.Sprintf("ALQ-%d", time.Now().UTC().UnixNano())
	return s.repo.Crear(ctx, alquiler)
}

func (s *AlquilerService) Actualizar(ctx context.Context, id int, empresaID int, alq *domain.Alquiler) (*domain.Alquiler, error) {
	old, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	alq.ID = id
	alq.EmpresaID = empresaID
	// Mantener el código original
	alq.Codigo = old.Codigo
	return s.repo.Actualizar(ctx, alq)
}

func (s *AlquilerService) Eliminar(ctx context.Context, id int, empresaID int) error {
	_, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	return s.repo.Eliminar(ctx, id)
}

func (s *AlquilerService) Terminar(ctx context.Context, id int, empresaID int) error {
	alq, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return err
	}
	alq.Estado = "finalizado"
	alq.ActivoParaCobro = false
	// Al actualizar a finalizado, el repo también debería liberar la unidad o podemos hacerlo explícito.
	// En mi implementación de repo.Eliminar lo hace, pero aquí es una actualización de estado.
	// Ajustaré el repo.Actualizar o lo haré aquí.
	_, err = s.repo.Actualizar(ctx, alq)
	return err
}

// --- Plantillas ---

func (s *AlquilerService) ListarPlantillas(ctx context.Context, empresaID int) ([]*domain.PlantillaContrato, error) {
	return s.repo.ListarPlantillas(ctx, empresaID)
}

func (s *AlquilerService) ObtenerPlantilla(ctx context.Context, id int, empresaID int) (*domain.PlantillaContrato, error) {
	return s.repo.ObtenerPlantilla(ctx, id, empresaID)
}

func (s *AlquilerService) GuardarPlantilla(ctx context.Context, p *domain.PlantillaContrato) (*domain.PlantillaContrato, error) {
	if p.ID > 0 {
		return s.repo.ActualizarPlantilla(ctx, p)
	}
	return s.repo.CrearPlantilla(ctx, p)
}

func (s *AlquilerService) EliminarPlantilla(ctx context.Context, id int, empresaID int) error {
	return s.repo.EliminarPlantilla(ctx, id, empresaID)
}

// --- Generación de Documentos ---

func (s *AlquilerService) GenerarContrato(ctx context.Context, id int, empresaID int, plantillaID int) (string, error) {
	alq, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return "", err
	}

	cliente, err := s.cliente.BuscarPorID(ctx, alq.ClienteID)
	if err != nil {
		return "", fmt.Errorf("no se pudo cargar datos del cliente")
	}

	var contenido string
	if plantillaID > 0 {
		plantilla, err := s.ObtenerPlantilla(ctx, plantillaID, empresaID)
		if err != nil {
			return "", fmt.Errorf("plantilla no encontrada")
		}
		contenido = plantilla.Contenido
	} else {
		// Plantilla por defecto básica
		contenido = `
# CONTRATO DE ALQUILER

Por el presente documento, se celebra un contrato de alquiler entre la empresa y el cliente **{{cliente_nombre}}**, identificado con documento **{{cliente_documento}}**.

## DETALLES DEL ALQUILER
- **Unidad:** {{unidad_codigo}}
- **Monto Mensual:** {{moneda}} {{monto_renta}}
- **Día de Pago:** {{dia_vencimiento}} de cada mes
- **Fecha de Inicio:** {{fecha_inicio}}

El arrendatario se compromete a cuidar la propiedad y realizar los pagos puntualmente.
`
	}

	// Reemplazar placeholders
	replacer := strings.NewReplacer(
		"{{cliente_nombre}}", alq.ClienteNombre,
		"{{cliente_documento}}", cliente.DocumentoNumero,
		"{{unidad_codigo}}", alq.UnidadCodigo,
		"{{monto_renta}}", fmt.Sprintf("%.2f", alq.MontoRenta),
		"{{moneda}}", alq.Moneda,
		"{{fecha_inicio}}", alq.FechaInicio.Format("2006-01-02"),
		"{{dia_vencimiento}}", fmt.Sprintf("%d", alq.DiaVencimiento),
	)

	return replacer.Replace(contenido), nil
}


type PagoAlquilerService struct {
	repo domain.PagoAlquilerRepository
}

func NewPagoAlquilerService(repo domain.PagoAlquilerRepository) *PagoAlquilerService {
	return &PagoAlquilerService{repo: repo}
}

func (s *PagoAlquilerService) Registrar(ctx context.Context, pago *domain.RegistroPagoAlquiler) (*domain.PagoAlquiler, error) {
	if pago.MesCorrespondiente < 1 || pago.MesCorrespondiente > 12 {
		return nil, fmt.Errorf("mes_correspondiente debe estar entre 1 y 12")
	}
	if pago.FechaPago.IsZero() {
		return nil, fmt.Errorf("fecha_pago es obligatoria")
	}
	return s.repo.Registrar(ctx, pago)
}

func (s *PagoAlquilerService) ListarPendientesMesActual(ctx context.Context, empresaID int) ([]*domain.PagoPendiente, error) {
	return s.repo.ListarPendientesMesActual(ctx, empresaID, time.Now().UTC())
}

func (s *PagoAlquilerService) ListarHistorial(ctx context.Context, filtros domain.PagoFiltros) ([]*domain.PagoAlquiler, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.Limite <= 0 {
		filtros.Limite = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *PagoAlquilerService) Obtener(ctx context.Context, id int, empresaID int) (*domain.PagoAlquiler, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Anular(ctx context.Context, id int, empresaID int) error {
	return s.repo.Eliminar(ctx, id, empresaID)
}

func (s *PagoAlquilerService) Actualizar(ctx context.Context, id int, empresaID int, notas *string, metodoPago string) (*domain.PagoAlquiler, error) {
	pago, err := s.Obtener(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	if notas != nil {
		pago.Nota = notas
	}
	if metodoPago != "" {
		pago.MetodoPago = metodoPago
	}
	return s.repo.Actualizar(ctx, pago)
}
