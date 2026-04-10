package domain

import "time"

type Cargo struct {
	ID                       int       `json:"id"`
	ContratoID               int       `json:"contrato_id"`
	Concepto                 string    `json:"concepto"`
	Descripcion              string    `json:"descripcion"`
	Moneda                   string    `json:"moneda"`
	PeriodoInicio            time.Time `json:"periodo_inicio"`
	PeriodoFin               time.Time `json:"periodo_fin"`
	FechaEmision             time.Time `json:"fecha_emision"`
	FechaVencimiento         time.Time `json:"fecha_vencimiento"`
	Monto                    float64   `json:"monto"`
	Saldo                    float64   `json:"saldo"`
	Estado                   string    `json:"estado"`
	GeneradoAutomaticamente  bool      `json:"generado_automaticamente"`
}

type CargoFiltros struct {
	EmpresaID  int
	ContratoID int
	Estado     string
	Pagina     int
	PorPagina  int
}

type RegistroCargo struct {
	ContratoID       int     `json:"contrato_id"`
	Concepto         string  `json:"concepto"`
	Descripcion      string  `json:"descripcion"`
	Monto            float64 `json:"monto"`
	FechaVencimiento string  `json:"fecha_vencimiento"`
}
