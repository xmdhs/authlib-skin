package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserToken holds the schema definition for the UserToken entity.
type UserToken struct {
	ent.Schema
}

// Fields of the UserToken.
func (UserToken) Fields() []ent.Field {
	return []ent.Field{
		// 用于验证 jwt token 是否被注销，若相同则有效
		field.Uint64("token_id"),
	}
}

// Edges of the UserToken.
func (UserToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("token").Unique(),
	}
}

func (UserToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user"),
	}
}
