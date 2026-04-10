package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Ticket struct {
	ent.Schema
}

func (Ticket) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Ticket) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tickets"},
	}
}

func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("unidad_id"),
		field.Int("cliente_id").Optional().Nillable(),
		field.String("asunto").MaxLen(120),
		field.String("descripcion").MaxLen(1000),
		field.Enum("prioridad").
			Values("baja", "media", "alta").
			Default("media"),
		field.Enum("estado").
			Values("abierto", "en_progreso", "resuelto", "cerrado").
			Default("abierto"),
	}
}

func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("tickets").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("unidad", Unidad.Type).
			Ref("tickets").
			Field("unidad_id").
			Required().
			Unique(),
		edge.From("cliente", Cliente.Type).
			Ref("tickets").
			Field("cliente_id").
			Unique(),
	}
}

func (Ticket) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id"),
		index.Fields("unidad_id"),
		index.Fields("estado"),
	}
}
