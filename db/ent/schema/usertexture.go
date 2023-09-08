package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserTexture holds the schema definition for the UserTexture entity.
type UserTexture struct {
	ent.Schema
}

// Fields of the UserTexture.
func (UserTexture) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_profile_id"),
		field.Int("texture_id"),
	}
}

// Edges of the UserTexture.
func (UserTexture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_profile", UserProfile.Type).
			Unique().
			Required().
			Field("user_profile_id"),
		edge.To("texture", Texture.Type).
			Unique().
			Required().
			Field("texture_id"),
	}
}

func (UserTexture) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user_profile"),
		index.Edges("texture"),
	}
}
