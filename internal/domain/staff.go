package domain

type Staff struct {
	ID        int     `json:"id"`
	UsuarioID int     `json:"usuario_id"`
	Usuario   string  `json:"usuario"`
	RolID     int     `json:"rol_id"`
	RolNombre string  `json:"rol_nombre"`
	EmpresaID int     `json:"empresa_id"`
	Principal bool    `json:"principal"`
	Estado    string  `json:"estado"`
}

type StaffFiltros struct {
	EmpresaID int
	Pagina    int
	PorPagina int
	Busqueda  string
}

type RegistroStaff struct {
	Usuario    string `json:"usuario"`
	Contrasena string `json:"contrasena"`
	RolID      int    `json:"rol_id"`
	EmpresaID  int    `json:"empresa_id"`
}
