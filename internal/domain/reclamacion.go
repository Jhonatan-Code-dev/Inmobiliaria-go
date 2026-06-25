package domain

import (
	"time"
)

type Reclamacion struct {
	ID                 int       `json:"id"`
	Codigo             string    `json:"codigo"`
	EmpresaID          int       `json:"empresa_id"`
	Nombres            string    `json:"nombres"`
	Apellidos          string    `json:"apellidos"`
	TipoDocumento      string    `json:"tipo_documento"`
	NumeroDocumento    string    `json:"numero_documento"`
	Telefono           string    `json:"telefono"`
	Email              string    `json:"email"`
	Direccion          string    `json:"direccion"`
	MenorEdad          bool      `json:"menor_edad"`
	NombreApoderado    string    `json:"nombre_apoderado,omitempty"`
	TipoBien           string    `json:"tipo_bien"` // PRODUCTO, SERVICIO
	MontoReclamado     float64   `json:"monto_reclamado"`
	DescripcionBien    string    `json:"descripcion_bien"`
	TipoReclamacion    string    `json:"tipo_reclamacion"` // RECLAMO, QUEJA
	DetalleReclamacion string    `json:"detalle_reclamacion"`
	PedidoConsumidor   string    `json:"pedido_consumidor"`
	Estado             string    `json:"estado"` // PENDIENTE, RESUELTO
	RespuestaDetalle   string    `json:"respuesta_detalle,omitempty"`
	RespondidoEn       *time.Time `json:"respondido_en,omitempty"`
	CreadoEn           time.Time `json:"creado_en"`
}

type ReclamacionFiltros struct {
	EmpresaID int
	Pagina    int
	Limite    int
}
