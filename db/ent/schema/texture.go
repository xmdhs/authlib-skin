package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Texture holds the schema definition for the Texture entity.
type Texture struct {
	ent.Schema
}

// Fields of the Texture.
func (Texture) Fields() []ent.Field {
	return []ent.Field{
		field.String("texture_hash").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(100)",
		}),
	}
}

// Edges of the Texture.
func (Texture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_user", User.Type).Unique().Required(),
		edge.To("user_profile", UserProfile.Type).Through("usertexture", UserTexture.Type),
	}
}

func (Texture) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("texture_hash").Unique(),
	}
}
