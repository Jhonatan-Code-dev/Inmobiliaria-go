package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type TipoPago struct {
	ent.Schema
}

func (TipoPago) Mixin() []ent.Mixin {
	return nil
}

func (TipoPago) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "tipos_pago"},
	}
}

func (TipoPago) Fields() []ent.Field {
	return []ent.Field{
		field.String("nombre").NotEmpty().Unique().MaxLen(50),
	}
}

func (TipoPago) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("gastos", Gasto.Type),
	}
}
