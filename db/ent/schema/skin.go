package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Skin holds the schema definition for the Skin entity.
type Skin struct {
	ent.Schema
}

// Fields of the Skin.
func (Skin) Fields() []ent.Field {
	return []ent.Field{
		field.String("skin_hash"),
		field.Uint8("type"),
		field.String("variant"),
	}
}

// Edges of the Skin.
func (Skin) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).Unique().Required(),
	}
}
