package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
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
		field.String("uuid").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(32)",
		}),
	}
}

// Edges of the UserToken.
func (UserToken) Edges() []ent.Edge {
	return nil
}

func (UserToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("uuid"),
	}
}
