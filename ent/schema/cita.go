package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Cita struct {
	ent.Schema
}

func (Cita) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Cita) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "citas"},
	}
}

func (Cita) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("propiedad_id").Optional().Nillable(),
		field.Int("unidad_id").Optional().Nillable(),
		field.Int("cliente_id").Optional().Nillable(),
		field.String("nombre_prospecto").NotEmpty().MaxLen(150),
		field.String("telefono_prospecto").NotEmpty().MaxLen(50),
		field.String("correo_prospecto").Optional().Nillable().MaxLen(150),
		field.Time("fecha_visita"),
		field.Enum("estado").
			Values("programada", "realizada", "cancelada", "no_asistio").
			Default("programada"),
		field.String("comentarios").Optional().Nillable().MaxLen(1000),
	}
}

func (Cita) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("citas").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("propiedad", Propiedad.Type).
			Ref("citas").
			Field("propiedad_id").
			Unique(),
		edge.From("unidad", Unidad.Type).
			Ref("citas").
			Field("unidad_id").
			Unique(),
		edge.From("cliente", Cliente.Type).
			Ref("citas").
			Field("cliente_id").
			Unique(),
	}
}

func (Cita) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id"),
		index.Fields("fecha_visita"),
	}
}
