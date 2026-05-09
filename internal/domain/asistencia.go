package domain

import "time"

// Horario define el turno de trabajo de un usuario
type Horario struct {
	ID                 int    `json:"id"`
	EmpresaID          int    `json:"empresa_id"`
	UsuarioID          int    `json:"usuario_id"`
	HoraEntrada        string `json:"hora_entrada"`       // Formato HH:mm
	HoraSalida         string `json:"hora_salida"`        // Formato HH:mm
	ToleranciaMinutos  int    `json:"tolerancia_minutos"` // Minutos de gracia
	DiasLaborables     string `json:"dias_laborables"`    // Ej: "1,2,3,4,5"
}

// RegistroHorario DTO para crear o actualizar un horario
type RegistroHorario struct {
	UsuarioID         int    `json:"usuario_id"`
	HoraEntrada       string `json:"hora_entrada"`
	HoraSalida        string `json:"hora_salida"`
	ToleranciaMinutos int    `json:"tolerancia_minutos"`
	DiasLaborables    string `json:"dias_laborables"`
}

// Asistencia representa una marca de entrada/salida de un trabajador
type Asistencia struct {
	ID              int        `json:"id"`
	EmpresaID       int        `json:"empresa_id"`
	UsuarioID       int        `json:"usuario_id"`
	UsuarioNombre   string     `json:"usuario_nombre"`
	Fecha           time.Time  `json:"fecha"`
	HoraEntrada     *time.Time `json:"hora_entrada"`
	HoraSalida      *time.Time `json:"hora_salida"`
	Estado          string     `json:"estado"` // puntual, tarde, falta, justificado
	Notas           *string    `json:"notas"`
	HorasTrabajadas *float64   `json:"horas_trabajadas"`
}

type AsistenciaFiltros struct {
	EmpresaID int
	UsuarioID int    // Opcional
	Estado    string // Opcional (puntual, tarde, etc)
	Desde     *time.Time
	Hasta     *time.Time
	Busqueda  string // Nueva: Búsqueda por nombre de trabajador
	Pagina    int
	Limite    int
}

// Permiso representa una justificación o solicitud de ausencia
type Permiso struct {
	ID            int       `json:"id"`
	EmpresaID     int       `json:"empresa_id"`
	UsuarioID     int       `json:"usuario_id"`
	UsuarioNombre string    `json:"usuario_nombre"`
	Fecha         time.Time `json:"fecha"`
	Motivo        string    `json:"motivo"`
	Estado        string    `json:"estado"` // pendiente, aprobado, rechazado
	Respuesta     *string   `json:"respuesta"`
}

// RegistroPermiso DTO para solicitar un permiso
type RegistroPermiso struct {
	Fecha  string `json:"fecha"`
	Motivo string `json:"motivo"`
}

// DecisionPermiso DTO para que el administrador apruebe/rechace
type DecisionPermiso struct {
	Estado    string `json:"estado"` // aprobado, rechazado
	Respuesta string `json:"respuesta"`
}

type PermisoFiltros struct {
	EmpresaID int
	UsuarioID int // Opcional
	Estado    string // Opcional
	Pagina    int
	Limite    int
}
