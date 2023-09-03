package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// UserToken holds the schema definition for the UserToken entity.
type UserToken struct {
	ent.Schema
}

// Fields of the UserToken.
func (UserToken) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("token_id"),
	}
}

// Edges of the UserToken.
func (UserToken) Edges() []ent.Edge {
	return nil
}
