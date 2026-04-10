package domain

import "time"

type Inmueble struct {
	ID             int        `json:"id"`
	EmpresaID      int        `json:"empresa_id"`
	Nombre         string     `json:"nombre"`
	Tipo           string     `json:"tipo"`
	Descripcion    *string    `json:"descripcion"`
	Direccion      string     `json:"direccion"`
	Ciudad         *string    `json:"ciudad"`
	Region         *string    `json:"region"`
	Pais           *string    `json:"pais"`
	CodigoPostal   *string    `json:"codigo_postal"`
	TotalPisos     int        `json:"total_pisos"`
	TotalUnidades  int        `json:"total_unidades"`
	Estado         string     `json:"estado"`
	CreadoEn       time.Time  `json:"creado_en"`
	Unidades       []*Unidad  `json:"unidades,omitempty"`
}

type Unidad struct {
	ID                 int        `json:"id"`
	PropiedadID        int        `json:"propiedad_id"`
	Codigo             string     `json:"codigo"`
	Nombre             *string    `json:"nombre"`
	Tipo               string     `json:"tipo"`
	NumeroPiso         *int       `json:"numero_piso"`
	Dormitorios        int        `json:"dormitorios"`
	Banos              int        `json:"banos"`
	AreaM2             *float64   `json:"area_m2"`
	Capacidad          int        `json:"capacidad"`
	Moneda             string     `json:"moneda"`
	PrecioBase         float64    `json:"precio_base"`
	PrecioBaseCents    int64      `json:"-"`
	DepositoRequerido  float64    `json:"deposito_requerido"`
	DepositoReqCents   int64      `json:"-"`
	IncluyeAgua        bool       `json:"incluye_agua"`
	IncluyeLuz         bool       `json:"incluye_luz"`
	IncluyeInternet    bool       `json:"incluye_internet"`
	Notas              *string    `json:"notas"`
	Estado             string     `json:"estado"`
	CreadoEn           time.Time  `json:"creado_en"`
}

type InmuebleFiltros struct {
	EmpresaID int
	Busqueda  string
	Pagina    int
	Limite    int
	Estado    string
	Tipo      string
}
