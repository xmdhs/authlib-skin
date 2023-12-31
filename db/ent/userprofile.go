// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
)

// UserProfile is the model entity for the UserProfile schema.
type UserProfile struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// UUID holds the value of the "uuid" field.
	UUID string `json:"uuid,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the UserProfileQuery when eager-loading is set.
	Edges        UserProfileEdges `json:"edges"`
	user_profile *int
	selectValues sql.SelectValues
}

// UserProfileEdges holds the relations/edges for other nodes in the graph.
type UserProfileEdges struct {
	// User holds the value of the user edge.
	User *User `json:"user,omitempty"`
	// Texture holds the value of the texture edge.
	Texture []*Texture `json:"texture,omitempty"`
	// Usertexture holds the value of the usertexture edge.
	Usertexture []*UserTexture `json:"usertexture,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// UserOrErr returns the User value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e UserProfileEdges) UserOrErr() (*User, error) {
	if e.loadedTypes[0] {
		if e.User == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: user.Label}
		}
		return e.User, nil
	}
	return nil, &NotLoadedError{edge: "user"}
}

// TextureOrErr returns the Texture value or an error if the edge
// was not loaded in eager-loading.
func (e UserProfileEdges) TextureOrErr() ([]*Texture, error) {
	if e.loadedTypes[1] {
		return e.Texture, nil
	}
	return nil, &NotLoadedError{edge: "texture"}
}

// UsertextureOrErr returns the Usertexture value or an error if the edge
// was not loaded in eager-loading.
func (e UserProfileEdges) UsertextureOrErr() ([]*UserTexture, error) {
	if e.loadedTypes[2] {
		return e.Usertexture, nil
	}
	return nil, &NotLoadedError{edge: "usertexture"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserProfile) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case userprofile.FieldID:
			values[i] = new(sql.NullInt64)
		case userprofile.FieldName, userprofile.FieldUUID:
			values[i] = new(sql.NullString)
		case userprofile.ForeignKeys[0]: // user_profile
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserProfile fields.
func (up *UserProfile) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case userprofile.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			up.ID = int(value.Int64)
		case userprofile.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				up.Name = value.String
			}
		case userprofile.FieldUUID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field uuid", values[i])
			} else if value.Valid {
				up.UUID = value.String
			}
		case userprofile.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field user_profile", value)
			} else if value.Valid {
				up.user_profile = new(int)
				*up.user_profile = int(value.Int64)
			}
		default:
			up.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the UserProfile.
// This includes values selected through modifiers, order, etc.
func (up *UserProfile) Value(name string) (ent.Value, error) {
	return up.selectValues.Get(name)
}

// QueryUser queries the "user" edge of the UserProfile entity.
func (up *UserProfile) QueryUser() *UserQuery {
	return NewUserProfileClient(up.config).QueryUser(up)
}

// QueryTexture queries the "texture" edge of the UserProfile entity.
func (up *UserProfile) QueryTexture() *TextureQuery {
	return NewUserProfileClient(up.config).QueryTexture(up)
}

// QueryUsertexture queries the "usertexture" edge of the UserProfile entity.
func (up *UserProfile) QueryUsertexture() *UserTextureQuery {
	return NewUserProfileClient(up.config).QueryUsertexture(up)
}

// Update returns a builder for updating this UserProfile.
// Note that you need to call UserProfile.Unwrap() before calling this method if this UserProfile
// was returned from a transaction, and the transaction was committed or rolled back.
func (up *UserProfile) Update() *UserProfileUpdateOne {
	return NewUserProfileClient(up.config).UpdateOne(up)
}

// Unwrap unwraps the UserProfile entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (up *UserProfile) Unwrap() *UserProfile {
	_tx, ok := up.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserProfile is not a transactional entity")
	}
	up.config.driver = _tx.drv
	return up
}

// String implements the fmt.Stringer.
func (up *UserProfile) String() string {
	var builder strings.Builder
	builder.WriteString("UserProfile(")
	builder.WriteString(fmt.Sprintf("id=%v, ", up.ID))
	builder.WriteString("name=")
	builder.WriteString(up.Name)
	builder.WriteString(", ")
	builder.WriteString("uuid=")
	builder.WriteString(up.UUID)
	builder.WriteByte(')')
	return builder.String()
}

// UserProfiles is a parsable slice of UserProfile.
type UserProfiles []*UserProfile
