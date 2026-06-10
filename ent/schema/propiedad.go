package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Propiedad struct {
	ent.Schema
}

func (Propiedad) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Propiedad) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "propiedades"},
	}
}

func (Propiedad) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.String("nombre").NotEmpty().MaxLen(150),
		field.Enum("tipo").
			Values("casa", "edificio", "quinta", "condominio", "otro").
			Default("casa"),
		field.String("descripcion").Optional().Nillable().MaxLen(1000),
		field.String("direccion").NotEmpty().MaxLen(255),
		field.String("ciudad").Optional().Nillable().MaxLen(120),
		field.String("region").Optional().Nillable().MaxLen(120),
		field.String("pais").Optional().Nillable().MaxLen(100),
		field.String("codigo_postal").Optional().Nillable().MaxLen(20),
		field.Int("total_pisos").Default(1).NonNegative(),
		field.Int("total_unidades").Default(1).NonNegative(),
		field.Enum("estado").
			Values("activa", "activo", "mantenimiento", "inactiva", "inactivo").
			Default("activa"),
	}
}

func (Propiedad) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("propiedades").
			Field("empresa_id").
			Required().
			Unique(),
		edge.To("unidades", Unidad.Type),
		edge.To("citas", Cita.Type),
	}
}
