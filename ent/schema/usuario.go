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
		entsql.Annotation{Table: "usuarios"},
	}
}

func (Usuario) Fields() []ent.Field {
	return []ent.Field{
		field.String("nombres").NotEmpty().MaxLen(120),
		field.String("apellidos").Optional().Nillable().MaxLen(120),
		field.String("correo").NotEmpty().Unique().MaxLen(150),
		field.String("telefono").Optional().Nillable().MaxLen(30),
		field.String("hash_contrasena").NotEmpty().Sensitive().MaxLen(255),
		field.Enum("estado").
			Values("activo", "inactivo", "bloqueado").
			Default("activo"),
		field.Time("ultimo_acceso").Optional().Nillable(),
	}
}

func (Usuario) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("empresas_usuario", EmpresaUsuario.Type),
	}
}
