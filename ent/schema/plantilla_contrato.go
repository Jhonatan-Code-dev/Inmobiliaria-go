package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type PlantillaContrato struct {
	ent.Schema
}

func (PlantillaContrato) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (PlantillaContrato) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "plantillas_contrato"},
	}
}

func (PlantillaContrato) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.String("nombre").NotEmpty().MaxLen(100),
		field.Text("contenido"),
	}
}

func (PlantillaContrato) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("plantillas_contrato").
			Field("empresa_id").
			Required().
			Unique(),
	}
}
