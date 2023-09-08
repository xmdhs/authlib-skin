// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
)

// Texture is the model entity for the Texture schema.
type Texture struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TextureHash holds the value of the "texture_hash" field.
	TextureHash string `json:"texture_hash,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Variant holds the value of the "variant" field.
	Variant string `json:"variant,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the TextureQuery when eager-loading is set.
	Edges                TextureEdges `json:"edges"`
	texture_created_user *int
	user_profile_texture *int
	selectValues         sql.SelectValues
}

// TextureEdges holds the relations/edges for other nodes in the graph.
type TextureEdges struct {
	// CreatedUser holds the value of the created_user edge.
	CreatedUser *User `json:"created_user,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// CreatedUserOrErr returns the CreatedUser value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e TextureEdges) CreatedUserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.CreatedUser == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.CreatedUser, nil
	}
	return nil, &NotLoadedError{edge: "created_user"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Texture) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case texture.FieldID:
			values[i] = new(sql.NullInt64)
		case texture.FieldTextureHash, texture.FieldType, texture.FieldVariant:
			values[i] = new(sql.NullString)
		case texture.ForeignKeys[0]: // texture_created_user
			values[i] = new(sql.NullInt64)
		case texture.ForeignKeys[1]: // user_profile_texture
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Texture fields.
func (t *Texture) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case texture.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			t.ID = int(value.Int64)
		case texture.FieldTextureHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field texture_hash", values[i])
			} else if value.Valid {
				t.TextureHash = value.String
			}
		case texture.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				t.Type = value.String
			}
		case texture.FieldVariant:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field variant", values[i])
			} else if value.Valid {
				t.Variant = value.String
			}
		case texture.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field texture_created_user", value)
			} else if value.Valid {
				t.texture_created_user = new(int)
				*t.texture_created_user = int(value.Int64)
			}
		case texture.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_profile_texture", value)
			} else if value.Valid {
				t.user_profile_texture = new(int)
				*t.user_profile_texture = int(value.Int64)
			}
		default:
			t.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Texture.
// This includes values selected through modifiers, order, etc.
func (t *Texture) Value(name string) (ent.Value, error) {
	return t.selectValues.Get(name)
}

// QueryCreatedUser queries the "created_user" edge of the Texture entity.
func (t *Texture) QueryCreatedUser() *UserQuery {
	return NewTextureClient(t.config).QueryCreatedUser(t)
}

// Update returns a builder for updating this Texture.
// Note that you need to call Texture.Unwrap() before calling this method if this Texture
// was returned from a transaction, and the transaction was committed or rolled back.
func (t *Texture) Update() *TextureUpdateOne {
	return NewTextureClient(t.config).UpdateOne(t)
}

// Unwrap unwraps the Texture entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (t *Texture) Unwrap() *Texture {
	_tx, ok := t.config.driver.(*txDriver)
	if !ok {
		panic("ent: Texture is not a transactional entity")
	}
	t.config.driver = _tx.drv
	return t
}

// String implements the fmt.Stringer.
func (t *Texture) String() string {
	var builder strings.Builder
	builder.WriteString("Texture(")
	builder.WriteString(fmt.Sprintf("id=%v, ", t.ID))
	builder.WriteString("texture_hash=")
	builder.WriteString(t.TextureHash)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(t.Type)
	builder.WriteString(", ")
	builder.WriteString("variant=")
	builder.WriteString(t.Variant)
	builder.WriteByte(')')
	return builder.String()
}

// Textures is a parsable slice of Texture.
type Textures []*Texture
