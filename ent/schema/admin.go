package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Admin representa credenciales para la consola administrativa
// y es independiente de los usuarios de empresas (multitenant).
type Admin struct {
	ent.Schema
}

func (Admin) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Admin) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "admins"},
	}
}

func (Admin) Fields() []ent.Field {
	return []ent.Field{
		field.String("nombre").NotEmpty().MaxLen(120),
		field.String("usuario").NotEmpty().Unique().MaxLen(80),
		field.String("hash_contrasena").NotEmpty().Sensitive().MaxLen(255),
		field.Bool("activo").Default(true),
	}
}
