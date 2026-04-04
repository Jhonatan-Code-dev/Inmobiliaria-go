package domain

type Usuario struct {
	ID             int
	Nombres        string
	Apellidos      string
	Correo         string
	Telefono       string
	HashContrasena string
	EmpresaID      int
}
