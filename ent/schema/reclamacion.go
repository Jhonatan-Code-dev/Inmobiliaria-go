package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Reclamacion struct {
	ent.Schema
}

func (Reclamacion) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

func (Reclamacion) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "reclamaciones"},
	}
}

func (Reclamacion) Fields() []ent.Field {
	return []ent.Field{
		field.String("codigo").NotEmpty().Unique().MaxLen(50),
		field.Int("empresa_id"),
		field.String("nombres").NotEmpty().MaxLen(150),
		field.String("apellidos").NotEmpty().MaxLen(150),
		field.String("tipo_documento").NotEmpty().MaxLen(50),
		field.String("numero_documento").NotEmpty().MaxLen(50),
		field.String("telefono").NotEmpty().MaxLen(50),
		field.String("email").NotEmpty().MaxLen(100),
		field.String("direccion").NotEmpty().MaxLen(255),
		field.Bool("menor_edad").Default(false),
		field.String("nombre_apoderado").Optional().MaxLen(200),
		field.String("tipo_bien").NotEmpty().MaxLen(20), // PRODUCTO, SERVICIO
		field.Float("monto_reclamado").Default(0.0),
		field.String("descripcion_bien").NotEmpty().MaxLen(1000),
		field.String("tipo_reclamacion").NotEmpty().MaxLen(20), // RECLAMO, QUEJA
		field.String("detalle_reclamacion").NotEmpty().MaxLen(4000),
		field.String("pedido_consumidor").NotEmpty().MaxLen(2000),
		field.String("estado").Default("PENDIENTE").MaxLen(20), // PENDIENTE, RESUELTO
		field.String("respuesta_detalle").Optional().MaxLen(4000),
		field.Time("respondido_en").Optional().Nillable(),
	}
}

func (Reclamacion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("reclamaciones").
			Field("empresa_id").
			Required().
			Unique(),
	}
}
