package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Asistencia struct {
	ent.Schema
}

func (Asistencia) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Asistencia) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asistencias"},
	}
}

func (Asistencia) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("usuario_id"),
		field.Time("fecha").Comment("Fecha lógica de la asistencia"),
		field.Time("hora_entrada").Optional().Nillable(),
		field.Time("hora_salida").Optional().Nillable(),
		field.Enum("estado").
			Values("puntual", "tarde", "falta", "justificado").
			Default("puntual"),
		field.String("notas").Optional().Nillable().MaxLen(500),
		field.Float("horas_trabajadas").Optional().Nillable(),
	}
}

func (Asistencia) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("asistencias").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("usuario", Usuario.Type).
			Ref("asistencias").
			Field("usuario_id").
			Required().
			Unique(),
	}
}

func (Asistencia) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "usuario_id", "fecha").Unique(),
	}
}
