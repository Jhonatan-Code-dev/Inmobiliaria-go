package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Empresa struct {
	ent.Schema
}

func (Empresa) Mixin() []ent.Mixin {
	return []ent.Mixin{AuditMixin{}}
}

func (Empresa) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "empresas",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
	}
}

func (Empresa) Fields() []ent.Field {
	return []ent.Field{
		field.String("nombre").NotEmpty().MaxLen(150),
		field.String("pais").Optional().Nillable().MaxLen(2),
		codigoMoneda("moneda", "PEN"),
		field.Int("maximo_usuarios").Default(1).Positive(),
		field.Bool("estado").Default(true),
		field.Time("vencimiento").Optional().Nillable(),
		field.String("horario_entrada_defecto").Default("08:00"),
		field.String("horario_salida_defecto").Default("17:00"),
		field.Int("tolerancia_defecto").Default(15),
		field.String("dias_laborables_defecto").Default("1,2,3,4,5"),
	}
}

func (Empresa) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("usuarios_empresa", EmpresaUsuario.Type),
		edge.To("clientes", Cliente.Type),
		edge.To("propiedades", Propiedad.Type),
		edge.To("contratos", Contrato.Type),
		edge.To("pagos", Pago.Type),
		edge.To("gastos", Gasto.Type),
		edge.To("movimientos_caja", MovimientoCaja.Type),
		edge.To("tickets", Ticket.Type),
		edge.To("horarios", Horario.Type),
		edge.To("asistencias", Asistencia.Type),
		edge.To("permisos", Permiso.Type),
		edge.To("plantillas_contrato", PlantillaContrato.Type),
		edge.To("citas", Cita.Type),
		edge.To("reclamaciones", Reclamacion.Type),
	}
}

func (Empresa) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("nombre"),
	}
}
