package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ServicioMedicion struct {
	ent.Schema
}

func (ServicioMedicion) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (ServicioMedicion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "servicio_mediciones"},
	}
}

func (ServicioMedicion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("unidad_id"),
		field.Enum("tipo_servicio").
			Values("agua", "luz", "gas", "internet", "otro").
			Default("agua"),
		fechaSolo("periodo_inicio"),
		fechaSolo("periodo_fin"),
		field.Float("lectura_anterior").Default(0),
		field.Float("lectura_actual").Default(0),
		field.Float("consumo").Default(0),
		codigoMoneda("moneda", "PEN"),
		montoExacto("tarifa_unitaria"),
		montoExacto("monto_total"),
		field.String("observaciones").Optional().Nillable().MaxLen(1000),
	}
}

func (ServicioMedicion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("unidad", Unidad.Type).
			Ref("servicio_mediciones").
			Field("unidad_id").
			Required().
			Unique(),
	}
}

func (ServicioMedicion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("unidad_id", "tipo_servicio", "periodo_inicio", "periodo_fin").Unique(),
	}
}
