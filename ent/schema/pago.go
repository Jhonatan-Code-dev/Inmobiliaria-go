package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Pago struct {
	ent.Schema
}

func (Pago) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Pago) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pagos"},
	}
}

func (Pago) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("cliente_id").Optional().Nillable(),
		field.Int("contrato_id").Optional().Nillable(),
		field.String("numero_recibo").NotEmpty().MaxLen(40),
		field.Time("fecha_pago"),
		codigoMoneda("moneda", "PEN"),
		montoExacto("monto_total"),
		field.Enum("metodo").
			Values("efectivo", "transferencia", "yape", "plin", "tarjeta", "deposito", "otro").
			Default("efectivo"),
		field.String("referencia").Optional().Nillable().MaxLen(120),
		field.String("notas").Optional().Nillable().MaxLen(1000),
		field.Enum("estado").
			Values("registrado", "confirmado", "anulado").
			Default("confirmado"),
	}
}

func (Pago) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("pagos").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("cliente", Cliente.Type).
			Ref("pagos").
			Field("cliente_id").
			Unique(),
		edge.From("contrato", Contrato.Type).
			Ref("pagos").
			Field("contrato_id").
			Unique(),
		edge.To("aplicaciones", PagoAplicacion.Type),
		edge.To("movimientos_caja", MovimientoCaja.Type),
	}
}

func (Pago) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "numero_recibo").Unique(),
	}
}
