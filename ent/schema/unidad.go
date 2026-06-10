package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Unidad struct {
	ent.Schema
}

func (Unidad) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Unidad) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "unidades"},
	}
}

func (Unidad) Fields() []ent.Field {
	return []ent.Field{
		field.Int("propiedad_id"),
		field.String("codigo").NotEmpty().MaxLen(30),
		field.String("nombre").Optional().Nillable().MaxLen(120),
		field.Enum("tipo").
			Values("cuarto", "departamento", "casa", "suite", "local", "otro").
			Default("cuarto"),
		field.Int("numero_piso").Optional().Nillable(),
		field.Int("dormitorios").Default(0).NonNegative(),
		field.Int("banos").Default(0).NonNegative(),
		field.Float("area_m2").Optional().Nillable(),
		field.Int("capacidad").Default(1).Positive(),
		codigoMoneda("moneda", "PEN"),
		montoExacto("precio_base"),
		montoExacto("deposito_requerido"),
		field.Bool("incluye_agua").Default(false),
		field.Bool("incluye_luz").Default(false),
		field.Bool("incluye_internet").Default(false),
		field.String("notas").Optional().Nillable().MaxLen(1000),
		field.Enum("estado").
			Values("disponible", "ocupado", "reservado", "mantenimiento", "inactiva").
			Default("disponible"),
	}
}

func (Unidad) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("propiedad", Propiedad.Type).
			Ref("unidades").
			Field("propiedad_id").
			Required().
			Unique(),
		edge.To("contratos", Contrato.Type),
		edge.To("servicio_mediciones", ServicioMedicion.Type),
		edge.To("tickets", Ticket.Type),
		edge.To("citas", Cita.Type),
	}
}

func (Unidad) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("propiedad_id", "codigo").Unique(),
	}
}
