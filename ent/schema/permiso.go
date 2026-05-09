package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Permiso struct {
	ent.Schema
}

func (Permiso) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Permiso) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "permisos"},
	}
}

func (Permiso) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("usuario_id"),
		field.Time("fecha").Comment("Fecha para la cual se solicita el permiso"),
		field.String("motivo").MaxLen(1000),
		field.Enum("estado").
			Values("pendiente", "aprobado", "rechazado").
			Default("pendiente"),
		field.String("respuesta").Optional().Nillable().MaxLen(1000).Comment("Respuesta del administrador"),
	}
}

func (Permiso) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("permisos").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("usuario", Usuario.Type).
			Ref("permisos").
			Field("usuario_id").
			Required().
			Unique(),
	}
}

func (Permiso) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "usuario_id"),
	}
}
