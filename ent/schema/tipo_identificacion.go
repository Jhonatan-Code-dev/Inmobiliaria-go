package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type TipoIdentificacion struct {
	ent.Schema
}

func (TipoIdentificacion) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (TipoIdentificacion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tipos_identificacion"},
	}
}

func (TipoIdentificacion) Fields() []ent.Field {
	return []ent.Field{
		field.String("codigo").NotEmpty().MaxLen(20),
		field.String("nombre").NotEmpty().MaxLen(80),
		field.String("pais").Optional().Nillable().MaxLen(2),
		field.Bool("activo").Default(true),
	}
}

func (TipoIdentificacion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("clientes", Cliente.Type),
	}
}

func (TipoIdentificacion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("codigo").Unique(),
	}
}
