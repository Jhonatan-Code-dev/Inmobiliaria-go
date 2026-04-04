package domain

type Empresa struct {
	ID              int
	Nombre          string
	DocumentoFiscal string
	Correo          string
	Telefono        string
	Direccion       string
	Ciudad          string
	Pais            string
	Moneda          string
	MaximoUsuarios  int
	Estado          string
}
