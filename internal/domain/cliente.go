package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrForbidden       = errors.New("forbidden")
	ErrClienteConDatos = errors.New("no se puede eliminar el cliente porque tiene registros asociados (contratos, pagos, etc)")
)

type Cliente struct {
	ID                   int        `json:"id"`
	EmpresaID            int        `json:"empresa_id"`
	TipoIdentificacionID int        `json:"tipo_identificacion_id"`
	DocumentoNumero      string     `json:"documento_numero"`
	Nombres              string     `json:"nombres"`
	Apellidos            *string    `json:"apellidos"`
	Correo               *string    `json:"correo"`
	FechaNacimiento      *time.Time `json:"fecha_nacimiento"`
	Nacionalidad         *string    `json:"nacionalidad"`
	Direccion            *string    `json:"direccion"`
	ContactoEmergencia   *string    `json:"contacto_emergencia"`
	TelefonoEmergencia   *string    `json:"telefono_emergencia"`
	Notas                *string    `json:"notas"`
	Estado               string     `json:"estado"`
	CreadoEn             time.Time  `json:"creado_en"`
}

type ClienteFiltros struct {
	EmpresaID int
	Busqueda  string // Para nombres, apellidos o documento
	Pagina    int
	Limite    int
}

type TipoIdentificacion struct {
	ID     int     `json:"id"`
	Codigo string  `json:"codigo"`
	Nombre string  `json:"nombre"`
	Pais   *string `json:"pais"`
	Activo bool    `json:"activo"`
}
