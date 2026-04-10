package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type MovimientoCaja struct {
	ent.Schema
}

func (MovimientoCaja) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (MovimientoCaja) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "movimientos_caja"},
	}
}

func (MovimientoCaja) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("pago_id").Optional().Nillable(),
		field.Int("gasto_id").Optional().Nillable(),
		field.Enum("tipo").
			Values("ingreso", "egreso", "ajuste").
			Default("ingreso"),
		field.String("concepto").NotEmpty().MaxLen(150),
		field.Time("fecha_movimiento"),
		codigoMoneda("moneda", "PEN"),
		field.Float("monto").
			SchemaType(map[string]string{
				dialect.MySQL: "decimal(12,2)",
			}).
			Default(0),
		field.Enum("metodo").
			Values("efectivo", "transferencia", "yape", "plin", "tarjeta", "deposito", "otro").
			Default("efectivo"),
		field.String("referencia").Optional().Nillable().MaxLen(120),
		field.String("observaciones").Optional().Nillable().MaxLen(1000),
	}
}

func (MovimientoCaja) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("movimientos_caja").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("pago", Pago.Type).
			Ref("movimientos_caja").
			Field("pago_id").
			Unique(),
		edge.From("gasto", Gasto.Type).
			Ref("movimientos_caja").
			Field("gasto_id").
			Unique(),
	}
}
