package schema

import (
	"rentals-go/internal/pkg/moneda"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

func codigoMoneda(nombre, defecto string) ent.Field {
	return field.String(nombre).
		Default(defecto).
		MaxLen(3).
		Validate(moneda.ValidarCodigo)
}

func montoExacto(nombre string) ent.Field {
	return field.Int64(nombre).
		Default(0)
}
