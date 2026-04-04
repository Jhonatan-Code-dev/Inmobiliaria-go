package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

func fechaSolo(nombre string) ent.Field {
	return field.Time(nombre).
		SchemaType(map[string]string{
			dialect.MySQL:    "date",
			dialect.Postgres: "date",
		})
}

func fechaSoloOpcional(nombre string) ent.Field {
	return field.Time(nombre).
		Optional().
		Nillable().
		SchemaType(map[string]string{
			dialect.MySQL:    "date",
			dialect.Postgres: "date",
		})
}
