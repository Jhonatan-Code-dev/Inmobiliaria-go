package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Horario struct {
	ent.Schema
}

func (Horario) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Horario) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "horarios"},
	}
}

func (Horario) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("usuario_id").Unique(), // Cada usuario tiene un único horario activo por ahora
		field.String("hora_entrada").Default("08:00").Comment("Formato HH:mm"),
		field.String("hora_salida").Default("17:00").Comment("Formato HH:mm"),
		field.Int("tolerancia_minutos").Default(15),
		field.String("dias_laborables").Default("1,2,3,4,5").Comment("Días de la semana ej: 1=Lunes, 5=Viernes"),
	}
}

func (Horario) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("horarios").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("usuario", Usuario.Type).
			Ref("horario").
			Field("usuario_id").
			Required().
			Unique(),
	}
}

func (Horario) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "usuario_id"),
	}
}
