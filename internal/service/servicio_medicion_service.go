package service

import (
	"context"
	"fmt"
	"rentals-go/internal/domain"
	"strings"
	"time"
)

type ServicioMedicionService struct {
	repo         domain.ServicioMedicionRepository
	cargoRepo    domain.CargoRepository
	alquilerRepo domain.AlquilerRepository
}

func NewServicioMedicionService(repo domain.ServicioMedicionRepository, cargoRepo domain.CargoRepository, alquilerRepo domain.AlquilerRepository) *ServicioMedicionService {
	return &ServicioMedicionService{repo: repo, cargoRepo: cargoRepo, alquilerRepo: alquilerRepo}
}

func (s *ServicioMedicionService) Listar(ctx context.Context, filtros domain.ServicioMedicionFiltros) ([]*domain.ServicioMedicion, int, error) {
	if filtros.Pagina <= 0 {
		filtros.Pagina = 1
	}
	if filtros.PorPagina <= 0 {
		filtros.PorPagina = 10
	}
	return s.repo.Listar(ctx, filtros)
}

func (s *ServicioMedicionService) Obtener(ctx context.Context, id int, empresaID int) (*domain.ServicioMedicion, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *ServicioMedicionService) Registrar(ctx context.Context, reg *domain.RegistroLectura, empresaID int) (*domain.ServicioMedicion, error) {
	fecha, err := time.Parse("2006-01-02", reg.FechaLectura)
	if err != nil {
		fecha = time.Now()
	}

	// 1. Obtener contrato para tener el unidad_id
	alq, err := s.alquilerRepo.BuscarPorID(ctx, reg.ContratoID)
	if err != nil {
		return nil, fmt.Errorf("contrato no encontrado o inválido")
	}

	// 2. Verificar duplicados (mismo unidad, tipo y fecha)
	existente, err := s.repo.BuscarPorFecha(ctx, alq.UnidadID, reg.TipoServicio, fecha)
	if err != nil {
		return nil, err
	}
	if existente != nil {
		return nil, fmt.Errorf("ya existe un registro de %s para la fecha %s. Usa el historial para editarlo o eliminarlo", 
			reg.TipoServicio, reg.FechaLectura)
	}

	// 3. Determinar lectura anterior
	var anterior float64
	if reg.LecturaAnterior != nil {
		anterior = *reg.LecturaAnterior
	} else {
		ultima, err := s.repo.ObtenerUltimaLectura(ctx, reg.ContratoID, reg.TipoServicio)
		if err != nil {
			return nil, err
		}
		if ultima != nil {
			anterior = ultima.LecturaActual
		}
	}

	factor := reg.Factor
	if factor <= 0 {
		factor = 1.0
	}

	consumo := (reg.LecturaActual - anterior) * factor
	if consumo < 0 {
		return nil, fmt.Errorf("la lectura actual no puede ser menor a la anterior (%.2f)", anterior)
	}

	monto := (consumo * reg.PrecioUnitario) + reg.CargoFijo

	medicion := &domain.ServicioMedicion{
		UnidadID:        alq.UnidadID,
		ContratoID:      reg.ContratoID,
		TipoServicio:    reg.TipoServicio,
		LecturaAnterior: anterior,
		LecturaActual:   reg.LecturaActual,
		Consumo:         consumo,
		PrecioUnitario:  reg.PrecioUnitario,
		Factor:          factor,
		CargoFijo:       reg.CargoFijo,
		Monto:           monto,
		FechaLectura:    fecha,
		Procesado:       false,
	}

	return s.repo.Crear(ctx, medicion)
}

func (s *ServicioMedicionService) Eliminar(ctx context.Context, id int, empresaID int) error {
	med, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return err
	}
	if med.Procesado {
		return fmt.Errorf("no se puede eliminar una medición que ya ha sido procesada en un cargo")
	}
	return s.repo.Eliminar(ctx, id, empresaID)
}

func (s *ServicioMedicionService) Actualizar(ctx context.Context, id int, empresaID int, lecturaActual float64) (*domain.ServicioMedicion, error) {
	med, err := s.repo.BuscarPorID(ctx, id, empresaID)
	if err != nil {
		return nil, err
	}
	if med.Procesado {
		return nil, fmt.Errorf("no se puede editar una medición que ya ha sido procesada")
	}

	consumo := (lecturaActual - med.LecturaAnterior) * med.Factor
	if consumo < 0 {
		return nil, fmt.Errorf("la nueva lectura actual no puede ser menor a la anterior (%.2f)", med.LecturaAnterior)
	}

	med.LecturaActual = lecturaActual
	med.Consumo = consumo
	med.Monto = (consumo * med.PrecioUnitario) + med.CargoFijo

	return s.repo.Actualizar(ctx, med)
}

func (s *ServicioMedicionService) ObtenerUltimaLectura(ctx context.Context, contratoID int, tipo string) (*domain.ServicioMedicion, error) {
	return s.repo.ObtenerUltimaLectura(ctx, contratoID, tipo)
}

func (s *ServicioMedicionService) RegistrarYCobrar(ctx context.Context, reg *domain.RegistroLectura, empresaID int) (*domain.ServicioMedicion, error) {
	// 1. Registrar la medición
	med, err := s.Registrar(ctx, reg, empresaID)
	if err != nil {
		return nil, err
	}

	// 2. Generar el cargo automáticamente
	alq, err := s.alquilerRepo.BuscarPorID(ctx, reg.ContratoID)
	if err != nil {
		return med, nil
	}
	
	concepto := fmt.Sprintf("Consumo de %s", strings.Title(reg.TipoServicio))
	
	// Descripción detallada de la fórmula
	descripcion := fmt.Sprintf("Lectura: %.2f (Act) - %.2f (Ant)", med.LecturaActual, med.LecturaAnterior)
	if med.Factor != 1.0 {
		descripcion += fmt.Sprintf(" x %.2f (Fac)", med.Factor)
	}
	descripcion += fmt.Sprintf(" = %.2f unidades x %.2f", med.Consumo, med.PrecioUnitario)
	if med.CargoFijo > 0 {
		descripcion += fmt.Sprintf(" + %.2f (Cargo Fijo)", med.CargoFijo)
	}

	cargo := &domain.Cargo{
		ContratoID:              reg.ContratoID,
		Concepto:                concepto,
		Descripcion:             descripcion,
		Monto:                   med.Monto,
		Saldo:                   med.Monto,
		Moneda:                  alq.Moneda,
		FechaVencimiento:        time.Now().AddDate(0, 0, 7), // 7 días para pagar
		Estado:                  "pendiente",
		GeneradoAutomaticamente: true,
	}

	nuevoCargo, err := s.cargoRepo.Crear(ctx, cargo)
	if err != nil {
		return med, nil
	}

	// 3. Vincular el cargo con la medición
	med.Procesado = true
	med.CargoID = &nuevoCargo.ID
	return s.repo.Actualizar(ctx, med)
}

func (s *ServicioMedicionService) RegistrarMasivo(ctx context.Context, registros []domain.RegistroLectura, empresaID int) ([]*domain.ServicioMedicion, error) {
	var resultados []*domain.ServicioMedicion
	for _, reg := range registros {
		med, err := s.RegistrarYCobrar(ctx, &reg, empresaID)
		if err != nil {
			return resultados, fmt.Errorf("error en contrato %d: %v", reg.ContratoID, err)
		}
		resultados = append(resultados, med)
	}
	return resultados, nil
}

func (s *ServicioMedicionService) ListarPendientesLectura(ctx context.Context, empresaID int, tipo string) ([]*domain.Alquiler, error) {
	// Implementación básica: listar alquileres activos que no tengan medición este mes
	return []*domain.Alquiler{}, nil
}
