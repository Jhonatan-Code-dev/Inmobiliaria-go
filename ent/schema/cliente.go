package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Cliente struct {
	ent.Schema
}

func (Cliente) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Cliente) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "clientes"},
	}
}

func (Cliente) Fields() []ent.Field {
	return []ent.Field{
		field.Int("empresa_id"),
		field.Int("tipo_identificacion_id"),
		field.String("documento_numero").NotEmpty().MaxLen(30),
		field.String("nombres").NotEmpty().MaxLen(120),
		field.String("apellidos").Optional().Nillable().MaxLen(120),
		field.String("correo").Optional().Nillable().MaxLen(150),
		fechaSoloOpcional("fecha_nacimiento"),
		field.String("nacionalidad").Optional().Nillable().MaxLen(60),
		field.String("direccion").Optional().Nillable().MaxLen(255),
		field.String("contacto_emergencia").Optional().Nillable().MaxLen(150),
		field.String("telefono_emergencia").Optional().Nillable().MaxLen(30),
		field.String("notas").Optional().Nillable().MaxLen(1000),
		field.Enum("estado").
			Values("activo", "inactivo", "bloqueado").
			Default("activo"),
	}
}

func (Cliente) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("empresa", Empresa.Type).
			Ref("clientes").
			Field("empresa_id").
			Required().
			Unique(),
		edge.From("tipo_identificacion", TipoIdentificacion.Type).
			Ref("clientes").
			Field("tipo_identificacion_id").
			Required().
			Unique(),
		edge.To("telefonos", ClienteTelefono.Type),
		edge.To("contratos", Contrato.Type),
		edge.To("pagos", Pago.Type),
		edge.To("tickets", Ticket.Type),
		edge.To("citas", Cita.Type),
	}
}

func (Cliente) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("empresa_id", "tipo_identificacion_id", "documento_numero").Unique(),
	}
}
