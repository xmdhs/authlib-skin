package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Unique().SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(30)",
		}),
		field.String("password").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(80)",
		}),
		field.String("salt").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(50)",
		}),
		field.String("reg_ip").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(32)",
		}),
		// 二进制状态位，保留
		field.Int("state"),
		field.Int64("reg_time"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("created_texture", Texture.Type).Ref("created_user"),
		edge.To("profile", UserProfile.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("token", UserToken.Type).Unique().Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("reg_ip"),
	}
}
