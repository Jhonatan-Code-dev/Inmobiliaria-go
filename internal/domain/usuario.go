package domain

type Usuario struct {
	ID             int
	Usuario        string
	HashContrasena string
	Estado         bool
	EmpresaID      int
}
