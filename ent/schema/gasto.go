package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Gasto struct {
	ent.Schema
}

func (Gasto) Mixin() []ent.Mixin {
	return nil
}

func (Gasto) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "gastos"},
	}
}

func (Gasto) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Float("monto").
			SchemaType(map[string]string{
				dialect.MySQL: "decimal(12,2)",
			}),
		field.Time("fecha"),
		field.Int("tipo_pago_id"),
		field.String("descripcion").NotEmpty().MaxLen(255),
	}
}

func (Gasto) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("gastos").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("tipo_pago", TipoPago.Type).
			Ref("gastos").
			Field("tipo_pago_id").
			Required().
			Unique(),
		edge.To("movimientos_caja", MovimientoCaja.Type),
	}
}
