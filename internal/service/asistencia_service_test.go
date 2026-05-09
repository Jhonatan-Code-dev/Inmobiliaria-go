package service

import (
	"context"
	"testing"
	"time"

	"rentals-go/internal/domain"
)

type horarioRepoStub struct {
	horario *domain.Horario
}

func (s *horarioRepoStub) BuscarPorUsuario(ctx context.Context, usuarioID int, empresaID int) (*domain.Horario, error) {
	return s.horario, nil
}
func (s *horarioRepoStub) Crear(ctx context.Context, h *domain.Horario) (*domain.Horario, error) {
	return h, nil
}
func (s *horarioRepoStub) Actualizar(ctx context.Context, h *domain.Horario) (*domain.Horario, error) {
	return h, nil
}

type asistenciaRepoStub struct {
	data    map[string]*domain.Asistencia
	created *domain.Asistencia
	updated *domain.Asistencia
}

func (s *asistenciaRepoStub) ListarPaginado(ctx context.Context, filtros domain.AsistenciaFiltros) ([]*domain.Asistencia, int, error) {
	return nil, 0, nil
}
func (s *asistenciaRepoStub) BuscarPorFechaUsuario(ctx context.Context, usuarioID int, empresaID int, fecha time.Time) (*domain.Asistencia, error) {
	key := fecha.Format("2006-01-02")
	return s.data[key], nil
}
func (s *asistenciaRepoStub) Crear(ctx context.Context, a *domain.Asistencia) (*domain.Asistencia, error) {
	s.created = a
	return a, nil
}
func (s *asistenciaRepoStub) Actualizar(ctx context.Context, a *domain.Asistencia) (*domain.Asistencia, error) {
	s.updated = a
	return a, nil
}
func (s *asistenciaRepoStub) Eliminar(ctx context.Context, id int, empresaID int) error {
	return nil
}

type permisoRepoStub struct{}

func (s *permisoRepoStub) ListarPaginado(ctx context.Context, filtros domain.PermisoFiltros) ([]*domain.Permiso, int, error) {
	return nil, 0, nil
}
func (s *permisoRepoStub) BuscarPorID(ctx context.Context, id int, empresaID int) (*domain.Permiso, error) {
	return nil, nil
}
func (s *permisoRepoStub) Crear(ctx context.Context, p *domain.Permiso) (*domain.Permiso, error) {
	return p, nil
}
func (s *permisoRepoStub) Actualizar(ctx context.Context, p *domain.Permiso) (*domain.Permiso, error) {
	return p, nil
}

type empresaRepoStub struct {
	empresa *domain.Empresa
}

func (s *empresaRepoStub) ListarPaginado(ctx context.Context, limite, offset int, busqueda string) ([]*domain.Empresa, int, error) {
	return nil, 0, nil
}
func (s *empresaRepoStub) BuscarPorID(ctx context.Context, id int) (*domain.Empresa, error) {
	return s.empresa, nil
}
func (s *empresaRepoStub) Crear(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	return emp, nil
}
func (s *empresaRepoStub) Actualizar(ctx context.Context, emp *domain.Empresa) (*domain.Empresa, error) {
	return emp, nil
}
func (s *empresaRepoStub) Eliminar(ctx context.Context, id int) error {
	return nil
}

func TestMarcarAsistencia_CalculoHoras(t *testing.T) {
	// Preparar datos
	hoy := time.Now().UTC()
	fechaHoy := time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, hoy.Location())
	
	// Simulamos que entró hace 8.5 horas
	entrada := hoy.Add(-8 * time.Hour).Add(-30 * time.Minute)
	
	repo := &asistenciaRepoStub{
		data: map[string]*domain.Asistencia{
			fechaHoy.Format("2006-01-02"): {
				ID:          1,
				EmpresaID:   1,
				UsuarioID:   1,
				Fecha:       fechaHoy,
				HoraEntrada: &entrada,
				Estado:      "puntual",
			},
		},
	}
	
	svc := NewAsistenciaService(&horarioRepoStub{}, repo, &permisoRepoStub{}, &empresaRepoStub{empresa: &domain.Empresa{Pais: "PE"}})
	
	// Ejecutar marcado de SALIDA
	res, err := svc.MarcarAsistencia(context.Background(), 1, 1)
	
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	
	if res.HoraSalida == nil {
		t.Fatal("se esperaba hora de salida")
	}
	
	if res.HorasTrabajadas == nil {
		t.Fatal("se esperaba calculo de horas trabajadas")
	}
	
	// Verificar calculo (8.5 horas aprox)
	esperado := 8.5
	if *res.HorasTrabajadas < esperado-0.01 || *res.HorasTrabajadas > esperado+0.01 {
		t.Errorf("calculo de horas incorrecto: obtuvo %v, se esperaba %v", *res.HorasTrabajadas, esperado)
	}
}

func TestMarcarAsistencia_DeteccionTardanza(t *testing.T) {
	// Horario: 08:00 con 15 min tolerancia
	horario := &domain.Horario{
		HoraEntrada:       "08:00",
		ToleranciaMinutos: 15,
	}
	
	repo := &asistenciaRepoStub{data: make(map[string]*domain.Asistencia)}
	hRepo := &horarioRepoStub{horario: horario}
	NewAsistenciaService(hRepo, repo, &permisoRepoStub{}, &empresaRepoStub{empresa: &domain.Empresa{Pais: "PE"}})
}
