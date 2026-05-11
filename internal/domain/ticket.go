package domain

import "time"

type Ticket struct {
	ID            int        `json:"id"`
	EmpresaID     int        `json:"empresa_id"`
	UnidadID      int        `json:"unidad_id"`
	UnidadNombre  string     `json:"unidad_nombre,omitempty"`
	ClienteID     *int       `json:"cliente_id,omitempty"`
	ClienteNombre string     `json:"cliente_nombre,omitempty"`
	Asunto        string     `json:"asunto"`
	Descripcion   string     `json:"descripcion"`
	Prioridad     string     `json:"prioridad"` // baja, media, alta
	Estado        string     `json:"estado"`    // abierto, en_progreso, resuelto, cerrado
	FechaApertura time.Time  `json:"fecha_apertura"`
	FechaCierre   *time.Time `json:"fecha_cierre,omitempty"`

}

type TicketFiltros struct {
	EmpresaID  int
	PropiedadID int
	UnidadID   int
	Estado    string
	Busqueda  string
	Pagina    int
	PorPagina int
}

type RegistroTicket struct {
	UnidadID    int    `json:"unidad_id"`
	ClienteID   *int   `json:"cliente_id,omitempty"`
	Asunto      string `json:"asunto"`
	Descripcion string `json:"descripcion"`
	Prioridad   string `json:"prioridad"`
}

type CambiarEstadoTicket struct {
	Estado string `json:"estado"` // abierto, en_progreso, resuelto, cerrado
}

// TicketResumen agrupa los tickets por estado para un vistazo rápido (Dashboard)
type TicketResumen struct {
	Total      int `json:"total"`
	Abiertos   int `json:"abiertos"`
	EnProgreso int `json:"en_progreso"`
	Resueltos  int `json:"resueltos"`
	Cerrados   int `json:"cerrados"`
}
