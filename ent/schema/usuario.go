package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Usuario struct {
	ent.Schema
}

func (Usuario) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Usuario) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "usuarios",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

func (Usuario) Fields() []ent.Field {
	return []ent.Field{
		field.String("usuario").NotEmpty().Unique().MaxLen(120),
		field.String("hash_contrasena").NotEmpty().Sensitive().MaxLen(255),
		field.Bool("estado").Default(true),
	}
}

func (Usuario) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("empresas_usuario", EmpresaUsuario.Type),
	}
}
