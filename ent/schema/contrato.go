package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Contrato struct {
	ent.Schema
}

func (Contrato) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Contrato) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "contratos"},
	}
}

func (Contrato) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("cliente_id"),
		field.Int("unidad_id"),
		field.String("codigo").NotEmpty().MaxLen(40),
		field.Enum("tipo").
			Values("alquiler", "reserva").
			Default("alquiler"),
		fechaSolo("fecha_inicio"),
		fechaSoloOpcional("fecha_fin"),
		field.Int("dia_vencimiento").Range(1, 32),
		codigoMoneda("moneda", "PEN"),
		montoExacto("monto_renta"),
		montoExacto("monto_deposito"),
		montoExacto("mora_diaria"),
		field.Bool("servicios_incluidos").Default(false),
		field.Bool("activo_para_cobro").Default(true),
		field.Enum("estado").
			Values("borrador", "activo", "vencido", "finalizado", "cancelado").
			Default("activo"),
		field.String("observaciones").Optional().Nillable().MaxLen(1500),
	}
}

func (Contrato) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("contratos").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("cliente", Cliente.Type).
			Ref("contratos").
			Field("cliente_id").
			Required().
			Unique(),
		edge.From("unidad", Unidad.Type).
			Ref("contratos").
			Field("unidad_id").
			Required().
			Unique(),
		edge.To("cargos", Cargo.Type),
		edge.To("pagos", Pago.Type),
		edge.To("servicio_mediciones", ServicioMedicion.Type),
	}
}

func (Contrato) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "codigo").Unique(),
	}
}
