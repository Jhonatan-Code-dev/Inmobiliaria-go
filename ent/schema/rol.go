package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Rol struct {
	ent.Schema
}

func (Rol) Mixin() []ent.Mixin {
	return nil
}

func (Rol) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "roles",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

func (Rol) Fields() []ent.Field {
	return []ent.Field{
		field.String("nombre").NotEmpty().Unique().MaxLen(60),
		field.String("descripcion").Optional().Nillable().MaxLen(255),
	}
}

func (Rol) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("usuarios_empresa", EmpresaUsuario.Type),
	}
}
