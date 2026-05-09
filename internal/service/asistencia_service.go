package service

import (
	"context"
	"fmt"
	"time"

	"rentals-go/internal/domain"
	"rentals-go/internal/pkg/tiempo"
)

type AsistenciaService struct {
	horarioRepo    domain.HorarioRepository
	asistenciaRepo domain.AsistenciaRepository
	permisoRepo    domain.PermisoRepository
	empresaRepo    domain.EmpresaRepository
}

func NewAsistenciaService(
	horarioRepo domain.HorarioRepository,
	asistenciaRepo domain.AsistenciaRepository,
	permisoRepo domain.PermisoRepository,
	empresaRepo domain.EmpresaRepository,
) *AsistenciaService {
	return &AsistenciaService{
		horarioRepo:    horarioRepo,
		asistenciaRepo: asistenciaRepo,
		permisoRepo:    permisoRepo,
		empresaRepo:    empresaRepo,
	}
}

// --- Horarios ---

func (s *AsistenciaService) ObtenerHorario(ctx context.Context, usuarioID int, empresaID int) (*domain.Horario, error) {
	return s.horarioRepo.BuscarPorUsuario(ctx, usuarioID, empresaID)
}

func (s *AsistenciaService) AsignarHorario(ctx context.Context, empresaID int, req *domain.RegistroHorario) (*domain.Horario, error) {
	// Verificar si ya existe
	existente, _ := s.horarioRepo.BuscarPorUsuario(ctx, req.UsuarioID, empresaID)

	horario := &domain.Horario{
		EmpresaID:         empresaID,
		UsuarioID:         req.UsuarioID,
		HoraEntrada:       req.HoraEntrada,
		HoraSalida:        req.HoraSalida,
		ToleranciaMinutos: req.ToleranciaMinutos,
		DiasLaborables:    req.DiasLaborables,
	}

	if existente != nil {
		horario.ID = existente.ID
		return s.horarioRepo.Actualizar(ctx, horario)
	}

	return s.horarioRepo.Crear(ctx, horario)
}

// --- Asistencia ---

func (s *AsistenciaService) MarcarAsistencia(ctx context.Context, usuarioID int, empresaID int) (*domain.Asistencia, error) {
	ahora := tiempo.AhoraUTC()
	// Idealmente se debería usar la zona horaria de la empresa, simplificamos con UTC para el demo
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, ahora.Location())

	// 1. Buscar si ya tiene una marca hoy
	registro, _ := s.asistenciaRepo.BuscarPorFechaUsuario(ctx, usuarioID, empresaID, hoy)

	// 2. Si no tiene registro, es ENTRADA
	if registro == nil {
		nuevoRegistro := &domain.Asistencia{
			EmpresaID:   empresaID,
			UsuarioID:   usuarioID,
			Fecha:       hoy,
			HoraEntrada: &ahora,
			Estado:      "puntual", // Default
		}

		// Obtener horario para verificar tardanza
		horario, err := s.horarioRepo.BuscarPorUsuario(ctx, usuarioID, empresaID)
		if err == nil && horario != nil {
			// Obtener zona horaria de la empresa
			zona := "UTC"
			emp, errEmp := s.empresaRepo.BuscarPorID(ctx, empresaID)
			if errEmp == nil && emp.Pais != "" {
				if emp.Pais == "PE" || emp.Pais == "CO" || emp.Pais == "EC" {
					zona = "America/Lima" // UTC-5
				} else if emp.Pais == "MX" {
					zona = "America/Mexico_City"
				} else if emp.Pais == "CL" {
					zona = "America/Santiago"
				}
			}

			// Convertir 'ahora' a la hora local de la empresa para comparar
			ahoraLocal, _ := tiempo.EnZona(ahora, zona)
			
			horaAsignada, errParse := time.Parse("15:04", horario.HoraEntrada)
			if errParse == nil {
				// Crear instancia de la hora esperada en la zona local
				horaEsperada := time.Date(ahoraLocal.Year(), ahoraLocal.Month(), ahoraLocal.Day(), horaAsignada.Hour(), horaAsignada.Minute(), 0, 0, ahoraLocal.Location())
				limiteTolerancia := horaEsperada.Add(time.Duration(horario.ToleranciaMinutos) * time.Minute)

				if ahoraLocal.After(limiteTolerancia) {
					nuevoRegistro.Estado = "tarde"
				}
			}
		}

		return s.asistenciaRepo.Crear(ctx, nuevoRegistro)
	}

	// 3. Si ya tiene registro, y no tiene salida, es SALIDA
	if registro.HoraSalida == nil {
		registro.HoraSalida = &ahora

		// Calcular horas trabajadas
		diff := ahora.Sub(*registro.HoraEntrada)
		horas := diff.Hours()
		registro.HorasTrabajadas = &horas

		return s.asistenciaRepo.Actualizar(ctx, registro)
	}

	// 4. Ya marcó entrada y salida
	return nil, fmt.Errorf("ya registró entrada y salida el día de hoy")
}

func (s *AsistenciaService) ListarAsistencia(ctx context.Context, filtros domain.AsistenciaFiltros) ([]*domain.Asistencia, int, error) {
	return s.asistenciaRepo.ListarPaginado(ctx, filtros)
}

func (s *AsistenciaService) ConsultarReporteAsistencia(ctx context.Context, filtros domain.AsistenciaFiltros) ([]*domain.Asistencia, int, error) {
	return s.asistenciaRepo.ConsultarReporteAsistencia(ctx, filtros)
}

func (s *AsistenciaService) ListarMiHistorial(ctx context.Context, usuarioID int, empresaID int) ([]*domain.Asistencia, error) {
	filtros := domain.AsistenciaFiltros{
		EmpresaID: empresaID,
		UsuarioID: usuarioID,
		Pagina:    1,
		Limite:    100, // Últimos 100 registros
	}
	asistencias, _, err := s.asistenciaRepo.ListarPaginado(ctx, filtros)
	return asistencias, err
}

// --- Permisos ---

func (s *AsistenciaService) SolicitarPermiso(ctx context.Context, usuarioID int, empresaID int, req *domain.RegistroPermiso) (*domain.Permiso, error) {
	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha inválido, use YYYY-MM-DD")
	}

	permiso := &domain.Permiso{
		EmpresaID: empresaID,
		UsuarioID: usuarioID,
		Fecha:     fecha,
		Motivo:    req.Motivo,
		Estado:    "pendiente",
	}

	return s.permisoRepo.Crear(ctx, permiso)
}

func (s *AsistenciaService) ListarPermisos(ctx context.Context, filtros domain.PermisoFiltros) ([]*domain.Permiso, int, error) {
	return s.permisoRepo.ListarPaginado(ctx, filtros)
}

func (s *AsistenciaService) DecidirPermiso(ctx context.Context, permisoID int, empresaID int, decision *domain.DecisionPermiso) (*domain.Permiso, error) {
	if decision.Estado != "aprobado" && decision.Estado != "rechazado" {
		return nil, fmt.Errorf("estado inválido")
	}

	existente, err := s.permisoRepo.BuscarPorID(ctx, permisoID, empresaID)
	if err != nil {
		return nil, err
	}

	existente.Estado = decision.Estado
	existente.Respuesta = &decision.Respuesta

	// Si es aprobado, podríamos crear un registro de asistencia como "justificado" o "permiso"
	// Para no complicar, solo actualizamos el permiso por ahora.

	return s.permisoRepo.Actualizar(ctx, existente)
}

func (s *AsistenciaService) EliminarAsistencia(ctx context.Context, id int, empresaID int) error {
	return s.asistenciaRepo.Eliminar(ctx, id, empresaID)
}
