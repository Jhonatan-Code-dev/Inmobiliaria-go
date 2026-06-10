package domain

import "time"

type Cita struct {
	ID                int        `json:"id"`
	EmpresaID         int        `json:"empresa_id"`
	PropiedadID       *int       `json:"propiedad_id,omitempty"`
	PropiedadNombre   string     `json:"propiedad_nombre,omitempty"`
	UnidadID          *int       `json:"unidad_id,omitempty"`
	UnidadNombre      string     `json:"unidad_nombre,omitempty"`
	ClienteID         *int       `json:"cliente_id,omitempty"`
	ClienteNombre     string     `json:"cliente_nombre,omitempty"`
	NombreProspecto   string     `json:"nombre_prospecto"`
	TelefonoProspecto string     `json:"telefono_prospecto"`
	CorreoProspecto   *string    `json:"correo_prospecto,omitempty"`
	FechaVisita       time.Time  `json:"fecha_visita"`
	Estado            string     `json:"estado"` // programada, realizada, cancelada, no_asistio
	Comentarios       *string    `json:"comentarios,omitempty"`
	CreadoEn          time.Time  `json:"creado_en"`
}

type CitaFiltros struct {
	EmpresaID   int
	PropiedadID int
	UnidadID    int
	Estado      string
	Busqueda    string
	Desde       *time.Time
	Hasta       *time.Time
	Pagina      int
	PorPagina   int
}

type RegistroCita struct {
	PropiedadID       *int      `json:"propiedad_id,omitempty"`
	UnidadID          *int      `json:"unidad_id,omitempty"`
	ClienteID         *int      `json:"cliente_id,omitempty"`
	NombreProspecto   string    `json:"nombre_prospecto"`
	TelefonoProspecto string    `json:"telefono_prospecto"`
	CorreoProspecto   *string   `json:"correo_prospecto,omitempty"`
	FechaVisita       time.Time `json:"fecha_visita"`
	Estado            string    `json:"estado"`
	Comentarios       *string   `json:"comentarios,omitempty"`
}
