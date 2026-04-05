package schema

import (
	"rentals-go/internal/pkg/tiempo"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type AuditMixin struct {
	mixin.Schema
}

func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("creado_en").
			Default(tiempo.AhoraUTC).
			Immutable(),
	}
}
