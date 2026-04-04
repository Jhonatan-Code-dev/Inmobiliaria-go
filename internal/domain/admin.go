package domain

// Admin representa un usuario de consola administrativa (single-tenant).
type Admin struct {
	ID              int
	Nombre          string
	Usuario         string
	HashContrasena  string
	Activo          bool
}
