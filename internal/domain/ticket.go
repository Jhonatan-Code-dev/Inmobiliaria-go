package domain

import "time"

type Ticket struct {
	ID            int       `json:"id"`
	EmpresaID     int       `json:"empresa_id"`
	UnidadID      int       `json:"unidad_id"`
	ClienteID     *int      `json:"cliente_id,omitempty"`
	Asunto        string    `json:"asunto"`
	Descripcion   string    `json:"descripcion"`
	Prioridad     string    `json:"prioridad"` // baja, media, alta
	Estado        string    `json:"estado"`    // abierto, en_progreso, cerrado
	FechaApertura time.Time `json:"fecha_apertura"`
	FechaCierre   *time.Time `json:"fecha_cierre,omitempty"`
}

type TicketFiltros struct {
	EmpresaID int
	UnidadID  int
	Estado    string
	Pagina    int
	PorPagina int
}

type RegistroTicket struct {
	UnidadID    int    `json:"unidad_id"`
	Asunto      string `json:"asunto"`
	Descripcion string `json:"descripcion"`
	Prioridad   string `json:"prioridad"`
}
