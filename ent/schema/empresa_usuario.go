package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type EmpresaUsuario struct {
	ent.Schema
}

func (EmpresaUsuario) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (EmpresaUsuario) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "empresa_usuarios"},
	}
}

func (EmpresaUsuario) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("usuario_id"),
		field.Int("rol_id"),
		field.Bool("principal").Default(false),
		field.Enum("estado").
			Values("activo", "invitado", "inactivo").
			Default("activo"),
	}
}

func (EmpresaUsuario) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("usuarios_empresa").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("usuario", Usuario.Type).
			Ref("empresas_usuario").
			Field("usuario_id").
			Required().
			Unique(),
		edge.From("rol", Rol.Type).
			Ref("usuarios_empresa").
			Field("rol_id").
			Required().
			Unique(),
	}
}

func (EmpresaUsuario) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "usuario_id").Unique(),
	}
}
