package domain

import "time"

type ServicioMedicion struct {
	ID                 int       `json:"id"`
	UnidadID           int       `json:"unidad_id"`
	ContratoID         int       `json:"contrato_id"`
	TipoServicio       string    `json:"tipo_servicio"` // agua, luz, etc.
	LecturaAnterior    float64   `json:"lectura_anterior"`
	LecturaActual      float64   `json:"lectura_actual"`
	Consumo            float64   `json:"consumo"`
	Monto              float64   `json:"monto"`
	FechaLectura       time.Time `json:"fecha_lectura"`
	Procesado          bool      `json:"procesado"`
	CargoID            *int      `json:"cargo_id,omitempty"`
}

type ServicioMedicionFiltros struct {
	EmpresaID  int
	ContratoID int
	Pagina     int
	PorPagina  int
}

type RegistroLectura struct {
	ContratoID    int     `json:"contrato_id"`
	TipoServicio  string  `json:"tipo_servicio"`
	LecturaActual float64 `json:"lectura_actual"`
	FechaLectura  string  `json:"fecha_lectura"`
	PrecioUnitario float64 `json:"precio_unitario"`
}
