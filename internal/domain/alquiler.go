package domain

import "time"

type Alquiler struct {
	ID               int        `json:"id"`
	EmpresaID        int        `json:"empresa_id"`
	ClienteID        int        `json:"cliente_id"`
	UnidadID         int        `json:"unidad_id"`
	Codigo           string     `json:"codigo"`
	Tipo             string     `json:"tipo"`
	FechaInicio      time.Time  `json:"fecha_inicio"`
	FechaFin         *time.Time `json:"fecha_fin"`
	DiaVencimiento   int        `json:"dia_vencimiento"`
	Moneda           string     `json:"moneda"`
	MontoRenta       float64    `json:"monto_renta"`
	MontoRentaCents  int64      `json:"-"`
	MontoDeposito    float64    `json:"monto_deposito"`
	MontoDepositoCts int64      `json:"-"`
	MoraDiaria       float64    `json:"mora_diaria"`
	MoraDiariaCents  int64      `json:"-"`
	ServiciosIncl    bool       `json:"servicios_incluidos"`
	ActivoParaCobro  bool       `json:"activo_para_cobro"`
	Estado           string     `json:"estado"`
	Observaciones    *string    `json:"observaciones"`
	CreadoEn         time.Time  `json:"creado_en"`

	ClienteNombre string `json:"cliente_nombre,omitempty"`
	UnidadCodigo  string `json:"unidad_codigo,omitempty"`
}

type AlquilerFiltros struct {
	EmpresaID int
	Busqueda  string
	Estado    string
	UnidadID  int
	Pagina    int
	Limite    int
}

type RegistroPagoAlquiler struct {
	EmpresaID          int
	ContratoID         int
	MontoPagado        float64
	MontoPagadoCents   int64
	FechaPago          time.Time
	MetodoPago         string
	Nota               *string
	MesCorrespondiente int
}

type PagoAlquiler struct {
	ID                 int       `json:"id"`
	EmpresaID          int       `json:"empresa_id"`
	ContratoID         int       `json:"alquiler_id"`
	ClienteID          *int      `json:"cliente_id"`
	NumeroRecibo       string    `json:"numero_recibo"`
	FechaPago          time.Time `json:"fecha_pago"`
	Moneda             string    `json:"moneda"`
	MontoPagado        float64   `json:"monto_pagado"`
	MontoPagadoCents   int64     `json:"-"`
	MetodoPago         string    `json:"metodo_pago"`
	Nota               *string   `json:"nota"`
	MesCorrespondiente int       `json:"mes_correspondiente"`
}

type PagoPendiente struct {
	AlquilerID      int       `json:"alquiler_id"`
	Cliente         string    `json:"cliente"`
	Unidad          string    `json:"unidad"`
	Monto           float64   `json:"monto"`
	MontoCents      int64     `json:"-"`
	FechaVencimiento time.Time `json:"fecha_vencimiento"`
	Estado          string    `json:"estado"`
}
