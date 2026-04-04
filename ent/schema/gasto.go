package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Gasto struct {
	ent.Schema
}

func (Gasto) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Gasto) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "gastos"},
	}
}

func (Gasto) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("propiedad_id").Optional().Nillable(),
		field.Int("unidad_id").Optional().Nillable(),
		field.Enum("categoria").
			Values("agua", "luz", "internet", "mantenimiento", "limpieza", "impuestos", "reparacion", "otro").
			Default("otro"),
		field.String("descripcion").NotEmpty().MaxLen(255),
		field.Time("fecha_gasto"),
		codigoMoneda("moneda", "PEN"),
		montoExacto("monto"),
		field.Enum("metodo_pago").
			Values("efectivo", "transferencia", "yape", "plin", "tarjeta", "deposito", "otro").
			Default("efectivo"),
		field.String("referencia").Optional().Nillable().MaxLen(120),
		field.String("pagado_a").Optional().Nillable().MaxLen(150),
		field.Enum("estado").
			Values("pendiente", "pagado", "anulado").
			Default("pagado"),
		field.String("notas").Optional().Nillable().MaxLen(1000),
	}
}

func (Gasto) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("gastos").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("propiedad", Propiedad.Type).
			Ref("gastos").
			Field("propiedad_id").
			Unique(),
		edge.From("unidad", Unidad.Type).
			Ref("gastos").
			Field("unidad_id").
			Unique(),
		edge.To("movimientos_caja", MovimientoCaja.Type),
	}
}
