package domain

import (
	"context"
	"time"
)

// ─────────────────────────────────────────
// Structs de KPIs / Dashboard
// ─────────────────────────────────────────

// ResumenGeneral contiene los KPIs principales del negocio en un vistazo.
type ResumenGeneral struct {
	// Inmuebles
	TotalPropiedades  int `json:"total_propiedades"`
	TotalUnidades     int `json:"total_unidades"`
	UnidadesOcupadas  int `json:"unidades_ocupadas"`
	UnidadesLibres    int `json:"unidades_libres"`
	TasaOcupacion     float64 `json:"tasa_ocupacion_pct"` // 0-100

	// Contratos
	ContratosActivos   int `json:"contratos_activos"`
	ContratosBorrador  int `json:"contratos_borrador"`
	ContratosVencidos  int `json:"contratos_vencidos"`

	// Finanzas del mes actual
	IngresosEsteMes    float64 `json:"ingresos_mes_actual"`
	GastosEsteMes      float64 `json:"gastos_mes_actual"`
	BalanceNeto        float64 `json:"balance_neto_mes"`

	// Morosidad
	TotalMorosos       int     `json:"total_morosos"`
	MontoPendiente     float64 `json:"monto_pendiente_cobro"`

	// Tickets
	TicketsAbiertos    int `json:"tickets_abiertos"`
	TicketsEnProgreso  int `json:"tickets_en_progreso"`

	// Clientes
	TotalClientes     int `json:"total_clientes"`

	// Periodo
	Mes  int `json:"mes"`
	Anio int `json:"anio"`
}

// OcupacionPropiedad muestra la ocupación desglosada por propiedad.
type OcupacionPropiedad struct {
	PropiedadID    int     `json:"propiedad_id"`
	Nombre         string  `json:"nombre"`
	Direccion      string  `json:"direccion"`
	TotalUnidades  int     `json:"total_unidades"`
	Ocupadas       int     `json:"ocupadas"`
	Libres         int     `json:"libres"`
	TasaOcupacion  float64 `json:"tasa_ocupacion_pct"`
}

// ResumenOcupacion contiene la ocupación global y por propiedad.
type ResumenOcupacion struct {
	TotalUnidades   int                  `json:"total_unidades"`
	TotalOcupadas   int                  `json:"total_ocupadas"`
	TotalLibres     int                  `json:"total_libres"`
	TasaGlobal      float64              `json:"tasa_global_pct"`
	PorPropiedad    []OcupacionPropiedad `json:"por_propiedad"`
}

// ClienteMoroso es un inquilino con pagos pendientes vencidos.
type ClienteMoroso struct {
	ClienteID       int       `json:"cliente_id"`
	NombreCompleto  string    `json:"nombre_completo"`
	UnidadCodigo    string    `json:"unidad_codigo"`
	PropiedadNombre string    `json:"propiedad_nombre"`
	ContratoID      int       `json:"contrato_id"`
	MontoPendiente  float64   `json:"monto_pendiente"`
	DiasVencido     int       `json:"dias_vencido"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
}

// ResumenMorosidad reúne todos los morosos y el total adeudado.
type ResumenMorosidad struct {
	TotalMorosos    int             `json:"total_morosos"`
	MontoTotal      float64         `json:"monto_total_pendiente"`
	Morosos         []ClienteMoroso `json:"morosos"`
}

// PuntoFinanciero representa un valor financiero en un periodo de tiempo.
type PuntoFinanciero struct {
	Periodo  string  `json:"periodo"`  // "2025-01", "2025-02" ...
	Ingresos float64 `json:"ingresos"`
	Gastos   float64 `json:"gastos"`
	Balance  float64 `json:"balance"`
}

// ReporteFinanciero contiene el flujo de ingresos y gastos en un rango de fechas.
type ReporteFinanciero struct {
	Desde          string            `json:"desde"`
	Hasta          string            `json:"hasta"`
	TotalIngresos  float64           `json:"total_ingresos"`
	TotalGastos    float64           `json:"total_gastos"`
	BalanceNeto    float64           `json:"balance_neto"`
	Serie          []PuntoFinanciero `json:"serie_mensual"`
}

// ContratoProximoVencer es un contrato próximo a terminar.
type ContratoProximoVencer struct {
	ContratoID     int       `json:"contrato_id"`
	Codigo         string    `json:"codigo"`
	ClienteNombre  string    `json:"cliente_nombre"`
	UnidadCodigo   string    `json:"unidad_codigo"`
	PropiedadNombre string   `json:"propiedad_nombre"`
	FechaFin       time.Time `json:"fecha_fin"`
	DiasRestantes  int       `json:"dias_restantes"`
	MontoRenta     float64   `json:"monto_renta"`
}

// EstadoCuentaCliente muestra el resumen financiero de un cliente.
type EstadoCuentaCliente struct {
	ClienteID      int     `json:"cliente_id"`
	NombreCompleto string  `json:"nombre_completo"`
	Documento      string  `json:"documento"`
	Correo         *string `json:"correo"`

	// Saldos
	TotalCargado   float64 `json:"total_cargado"`
	TotalPagado    float64 `json:"total_pagado"`
	SaldoPendiente float64 `json:"saldo_pendiente"`

	// Desglose de cargos
	Cargos []CargoResumen `json:"cargos"`
}

// CargoResumen es la versión liviana de un cargo para el estado de cuenta.
type CargoResumen struct {
	CargoID          int       `json:"cargo_id"`
	Concepto         string    `json:"concepto"`
	Monto            float64   `json:"monto"`
	Saldo            float64   `json:"saldo"`
	Estado           string    `json:"estado"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
}

// TopUnidad muestra las unidades que más ingresos generan.
type TopUnidad struct {
	UnidadID       int     `json:"unidad_id"`
	Codigo         string  `json:"codigo"`
	PropiedadNombre string  `json:"propiedad_nombre"`
	TotalIngresos  float64 `json:"total_ingresos"`
	TotalPagos     int     `json:"total_pagos"`
}

// DashboardFiltros parámetros de consulta reutilizables.
type DashboardFiltros struct {
	EmpresaID int
	Desde     time.Time
	Hasta     time.Time
	DiasAlerta int // para contratos por vencer
	ClienteID  int
}

// ─────────────────────────────────────────
// Puerto del servicio (interfaz)
// ─────────────────────────────────────────

// DashboardRepository define las consultas de base de datos necesarias.
type DashboardRepository interface {
	ResumenGeneral(ctx context.Context, empresaID int, ahora time.Time) (*ResumenGeneral, error)
	ResumenOcupacion(ctx context.Context, empresaID int) (*ResumenOcupacion, error)
	ResumenMorosidad(ctx context.Context, empresaID int, ahora time.Time) (*ResumenMorosidad, error)
	ReporteFinanciero(ctx context.Context, empresaID int, desde, hasta time.Time) (*ReporteFinanciero, error)
	ContratosProximosVencer(ctx context.Context, empresaID int, dias int, ahora time.Time) ([]ContratoProximoVencer, error)
	EstadoCuentaCliente(ctx context.Context, empresaID, clienteID int) (*EstadoCuentaCliente, error)
	TopUnidades(ctx context.Context, empresaID int, desde, hasta time.Time, limite int) ([]TopUnidad, error)
}

// DashboardService define la lógica de negocio del dashboard.
type DashboardService interface {
	ObtenerResumenGeneral(ctx context.Context, empresaID int) (*ResumenGeneral, error)
	ObtenerOcupacion(ctx context.Context, empresaID int) (*ResumenOcupacion, error)
	ObtenerMorosidad(ctx context.Context, empresaID int) (*ResumenMorosidad, error)
	ObtenerReporteFinanciero(ctx context.Context, empresaID int, desde, hasta time.Time) (*ReporteFinanciero, error)
	ObtenerContratosProximosVencer(ctx context.Context, empresaID int, dias int) ([]ContratoProximoVencer, error)
	ObtenerEstadoCuentaCliente(ctx context.Context, empresaID, clienteID int) (*EstadoCuentaCliente, error)
	ObtenerTopUnidades(ctx context.Context, empresaID int, desde, hasta time.Time, limite int) ([]TopUnidad, error)
}
