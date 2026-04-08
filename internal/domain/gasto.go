package domain

import (
	"time"
)

type Gasto struct {
	ID          int
	EmpresaID   int
	Monto       float64
	Fecha       time.Time
	TipoPagoID  int
	Descripcion string
}

type GastoFiltros struct {
	EmpresaID int
	Anio      int
	Mes       int
	Desde     *time.Time
	Hasta     *time.Time
	Fecha     *time.Time
	Pagina    int
	Limite    int
}

type TipoPago struct {
	ID     int
	Nombre string
}

type MovimientoCaja struct {
	ID              int
	EmpresaID       int
	PagoID          *int
	GastoID         *int
	Tipo            string // ingreso, egreso, ajuste
	Concepto        string
	FechaMovimiento time.Time
	Moneda          string
	Monto           int64
	Metodo          string
	Referencia      *string
	Observaciones   *string
}
