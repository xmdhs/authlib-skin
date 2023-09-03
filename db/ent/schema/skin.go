package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		edge.To("created_user", User.Type).Unique().Required(),
	}
}

func (Skin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("skin_hash"),
	}
}
