package domain

import (
	"time"
)

// PuntoIngresoGasto representa la serie mensual de ingresos y gastos
type PuntoIngresoGasto struct {
	Periodo  string  `json:"periodo"` // Formato "YYYY-MM"
	Ingresos float64 `json:"ingresos"`
	Gastos   float64 `json:"gastos"`
	Balance  float64 `json:"balance"`
}

// ReporteIngresosGastos contiene la serie total de balance financiero
type ReporteIngresosGastos struct {
	Desde         string              `json:"desde"`
	Hasta         string              `json:"hasta"`
	TotalIngresos float64             `json:"total_ingresos"`
	TotalGastos   float64             `json:"total_gastos"`
	BalanceNeto   float64             `json:"balance_neto"`
	Serie         []PuntoIngresoGasto `json:"serie_mensual"`
}

// DistribucionMetodoPago detalla el volumen de ingresos según el método usado
type DistribucionMetodoPago struct {
	Metodo        string  `json:"metodo"` // "efectivo", "transferencia", "yape", etc.
	Total         float64 `json:"total"`
	CantidadPagos int     `json:"cantidad_pagos"`
	Porcentaje    float64 `json:"porcentaje"`
}

// DistribucionCategoriaGasto detalla el volumen de gastos según su tipo/categoría
type DistribucionCategoriaGasto struct {
	TipoPagoID     int     `json:"tipo_pago_id"`
	Categoria      string  `json:"categoria"` // Nombre del TipoPago
	Total          float64 `json:"total"`
	CantidadGastos int     `json:"cantidad_gastos"`
	Porcentaje     float64 `json:"porcentaje"`
}

// RentabilidadPropiedad calcula la rentabilidad neta y ocupación por propiedad
type RentabilidadPropiedad struct {
	PropiedadID      int     `json:"propiedad_id"`
	Nombre           string  `json:"nombre"`
	Direccion        string  `json:"direccion"`
	TotalUnidades    int     `json:"total_unidades"`
	UnidadesOcupadas int     `json:"unidades_ocupadas"`
	TasaOcupacionPct float64 `json:"tasa_ocupacion_pct"`
	Ingresos         float64 `json:"ingresos"`
	Gastos           float64 `json:"gastos"`
	Rentabilidad     float64 `json:"rentabilidad"` // Ingresos - Gastos
}

// TicketsPorEstado cuenta los tickets según su estado actual
type TicketsPorEstado struct {
	Abierto    int `json:"abierto"`
	EnProgreso int `json:"en_progreso"`
	Resuelto   int `json:"resuelto"`
	Anulado    int `json:"anulado"`
}

// TicketsPorPrioridad cuenta los tickets según su nivel de urgencia
type TicketsPorPrioridad struct {
	Baja  int `json:"baja"`
	Media int `json:"media"`
	Alta  int `json:"alta"`
}

// ResumenMantenimientoReporte detalla métricas del módulo de soporte/tickets
type ResumenMantenimientoReporte struct {
	TotalTickets int                 `json:"total_tickets"`
	PorEstado    TicketsPorEstado    `json:"por_estado"`
	PorPrioridad TicketsPorPrioridad `json:"por_prioridad"`
}

// ReporteFiltros contiene parámetros comunes para reportes
type ReporteFiltros struct {
	EmpresaID int
	Desde     time.Time
	Hasta     time.Time
}
