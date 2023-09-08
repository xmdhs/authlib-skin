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
		// 皮肤 or 披风
		field.String("type").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(10)",
		}),
		// slim or 空
		field.String("variant").SchemaType(map[string]string{
			dialect.MySQL: "VARCHAR(10)",
		}),
	}
}

// Edges of the Texture.
func (Texture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("created_user", User.Type).Unique().Required(),
	}
}

func (Texture) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("texture_hash"),
	}
}
