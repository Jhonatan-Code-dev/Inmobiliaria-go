package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type PagoAplicacion struct {
	ent.Schema
}

func (PagoAplicacion) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (PagoAplicacion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pago_aplicaciones"},
	}
}

func (PagoAplicacion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("pago_id"),
		field.Int("cargo_id"),
		codigoMoneda("moneda", "PEN"),
		montoExacto("monto_aplicado"),
	}
}

func (PagoAplicacion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pago", Pago.Type).
			Ref("aplicaciones").
			Field("pago_id").
			Required().
			Unique(),
		edge.From("cargo", Cargo.Type).
			Ref("aplicaciones_pago").
			Field("cargo_id").
			Required().
			Unique(),
	}
}

func (PagoAplicacion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("pago_id", "cargo_id").Unique(),
	}
}
