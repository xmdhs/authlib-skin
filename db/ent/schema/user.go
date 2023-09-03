package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Unique(),
		field.String("password"),
		field.String("salt"),
		field.Int("state"),
		field.Int64("reg_time"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("skin", Skin.Type).Ref("user").Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("profile", UserProfile.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
