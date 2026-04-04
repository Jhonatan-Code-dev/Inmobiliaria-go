package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ClienteTelefono struct {
	ent.Schema
}

func (ClienteTelefono) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (ClienteTelefono) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cliente_telefonos"},
	}
}

func (ClienteTelefono) Fields() []ent.Field {
	return []ent.Field{
		field.Int("cliente_id"),
		field.String("telefono").NotEmpty().MaxLen(30),
		field.String("etiqueta").Optional().Nillable().MaxLen(30),
		field.Bool("principal").Default(false),
		field.Bool("whatsapp").Default(false),
	}
}

func (ClienteTelefono) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cliente", Cliente.Type).
			Ref("telefonos").
			Field("cliente_id").
			Required().
			Unique(),
	}
}
