package domain

import (
	"time"
)

type Empresa struct {
	ID              int
	Nombre          string
	Pais            string
	Moneda          string
	MaximoUsuarios  int
	Estado          bool
	Vencimiento     time.Time
	CreadoEn        time.Time
}
