package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserProfile holds the schema definition for the UserProfile entity.
type UserProfile struct {
	ent.Schema
}

// Fields of the UserProfile.
func (UserProfile) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(20)",
		}),
		field.String("uuid").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(32)",
		}),
	}
}

// Edges of the UserProfile.
func (UserProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("profile").Required().Unique(),
		edge.From("texture", Texture.Type).Ref("user_profile").Through("usertexture", UserTexture.Type),
	}
}

func (UserProfile) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user"),
		index.Fields("name"),
		index.Fields("uuid"),
	}
}
