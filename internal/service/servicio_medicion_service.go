package service

import (
	"context"
	"fmt"
	"rentals-go/internal/domain"
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
	
	ultima, err := s.repo.ObtenerUltimaLectura(ctx, reg.ContratoID, reg.TipoServicio)
	if err != nil {
		return nil, err
	}

	anterior := 0.0
	if ultima != nil {
		anterior = ultima.LecturaActual
	}

	consumo := reg.LecturaActual - anterior
	if consumo < 0 {
		return nil, fmt.Errorf("la lectura actual no puede ser menor a la anterior (%.2f)", anterior)
	}

	monto := consumo * reg.PrecioUnitario

	// Necesitamos el unidad_id del contrato
	alq, err := s.alquilerRepo.BuscarPorID(ctx, reg.ContratoID)
	if err != nil {
		return nil, fmt.Errorf("alquiler no encontrado")
	}

	medicion := &domain.ServicioMedicion{
		UnidadID:        alq.UnidadID,
		ContratoID:      reg.ContratoID,
		TipoServicio:    reg.TipoServicio,
		LecturaAnterior: anterior,
		LecturaActual:   reg.LecturaActual,
		Consumo:         consumo,
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

	consumo := lecturaActual - med.LecturaAnterior
	if consumo < 0 {
		return nil, fmt.Errorf("la nueva lectura actual no puede ser menor a la anterior (%.2f)", med.LecturaAnterior)
	}

	// Recalcular monto basado en la tarifa original (estimada del monto/consumo previo para no pedir precio unitario de nuevo)
	precioUnitario := 0.0
	if med.Consumo > 0 {
		precioUnitario = med.Monto / med.Consumo
	}
	
	med.LecturaActual = lecturaActual
	med.Consumo = consumo
	med.Monto = consumo * precioUnitario

	return s.repo.Actualizar(ctx, med)
}
