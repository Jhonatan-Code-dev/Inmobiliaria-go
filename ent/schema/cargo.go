package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Cargo struct {
	ent.Schema
}

func (Cargo) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Cargo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cargos"},
	}
}

func (Cargo) Fields() []ent.Field {
	return []ent.Field{
		field.Int("contrato_id"),
		field.Enum("concepto").
			Values("renta", "deposito", "agua", "luz", "internet", "mantenimiento", "mora", "otro").
			Default("renta"),
		field.String("descripcion").Optional().Nillable().MaxLen(255),
		codigoMoneda("moneda", "PEN"),
		fechaSolo("periodo_inicio"),
		fechaSolo("periodo_fin"),
		fechaSolo("fecha_emision"),
		fechaSolo("fecha_vencimiento"),
		montoExacto("monto"),
		montoExacto("saldo"),
		field.Enum("estado").
			Values("pendiente", "parcial", "pagado", "vencido", "anulado").
			Default("pendiente"),
		field.Bool("generado_automaticamente").Default(false),
	}
}

func (Cargo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("contrato", Contrato.Type).
			Ref("cargos").
			Field("contrato_id").
			Required().
			Unique(),
		edge.To("aplicaciones_pago", PagoAplicacion.Type),
	}
}

func (Cargo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("contrato_id", "fecha_vencimiento"),
	}
}
