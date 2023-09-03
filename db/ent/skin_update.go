// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/xmdhs/authlib-skin/db/ent/predicate"
	"github.com/xmdhs/authlib-skin/db/ent/skin"
	"github.com/xmdhs/authlib-skin/db/ent/user"
)

// SkinUpdate is the builder for updating Skin entities.
type SkinUpdate struct {
	config
	hooks    []Hook
	mutation *SkinMutation
}

// Where appends a list predicates to the SkinUpdate builder.
func (su *SkinUpdate) Where(ps ...predicate.Skin) *SkinUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetSkinHash sets the "skin_hash" field.
func (su *SkinUpdate) SetSkinHash(s string) *SkinUpdate {
	su.mutation.SetSkinHash(s)
	return su
}

// SetType sets the "type" field.
func (su *SkinUpdate) SetType(u uint8) *SkinUpdate {
	su.mutation.ResetType()
	su.mutation.SetType(u)
	return su
}

// AddType adds u to the "type" field.
func (su *SkinUpdate) AddType(u int8) *SkinUpdate {
	su.mutation.AddType(u)
	return su
}

// SetVariant sets the "variant" field.
func (su *SkinUpdate) SetVariant(s string) *SkinUpdate {
	su.mutation.SetVariant(s)
	return su
}

// SetUserID sets the "user" edge to the User entity by ID.
func (su *SkinUpdate) SetUserID(id int) *SkinUpdate {
	su.mutation.SetUserID(id)
	return su
}

// SetUser sets the "user" edge to the User entity.
func (su *SkinUpdate) SetUser(u *User) *SkinUpdate {
	return su.SetUserID(u.ID)
}

// Mutation returns the SkinMutation object of the builder.
func (su *SkinUpdate) Mutation() *SkinMutation {
	return su.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (su *SkinUpdate) ClearUser() *SkinUpdate {
	su.mutation.ClearUser()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SkinUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SkinUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SkinUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SkinUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SkinUpdate) check() error {
	if _, ok := su.mutation.UserID(); su.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Skin.user"`)
	}
	return nil
}

func (su *SkinUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(skin.Table, skin.Columns, sqlgraph.NewFieldSpec(skin.FieldID, field.TypeInt))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.SkinHash(); ok {
		_spec.SetField(skin.FieldSkinHash, field.TypeString, value)
	}
	if value, ok := su.mutation.GetType(); ok {
		_spec.SetField(skin.FieldType, field.TypeUint8, value)
	}
	if value, ok := su.mutation.AddedType(); ok {
		_spec.AddField(skin.FieldType, field.TypeUint8, value)
	}
	if value, ok := su.mutation.Variant(); ok {
		_spec.SetField(skin.FieldVariant, field.TypeString, value)
	}
	if su.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   skin.UserTable,
			Columns: []string{skin.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   skin.UserTable,
			Columns: []string{skin.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{skin.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SkinUpdateOne is the builder for updating a single Skin entity.
type SkinUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SkinMutation
}

// SetSkinHash sets the "skin_hash" field.
func (suo *SkinUpdateOne) SetSkinHash(s string) *SkinUpdateOne {
	suo.mutation.SetSkinHash(s)
	return suo
}

// SetType sets the "type" field.
func (suo *SkinUpdateOne) SetType(u uint8) *SkinUpdateOne {
	suo.mutation.ResetType()
	suo.mutation.SetType(u)
	return suo
}

// AddType adds u to the "type" field.
func (suo *SkinUpdateOne) AddType(u int8) *SkinUpdateOne {
	suo.mutation.AddType(u)
	return suo
}

// SetVariant sets the "variant" field.
func (suo *SkinUpdateOne) SetVariant(s string) *SkinUpdateOne {
	suo.mutation.SetVariant(s)
	return suo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (suo *SkinUpdateOne) SetUserID(id int) *SkinUpdateOne {
	suo.mutation.SetUserID(id)
	return suo
}

// SetUser sets the "user" edge to the User entity.
func (suo *SkinUpdateOne) SetUser(u *User) *SkinUpdateOne {
	return suo.SetUserID(u.ID)
}

// Mutation returns the SkinMutation object of the builder.
func (suo *SkinUpdateOne) Mutation() *SkinMutation {
	return suo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (suo *SkinUpdateOne) ClearUser() *SkinUpdateOne {
	suo.mutation.ClearUser()
	return suo
}

// Where appends a list predicates to the SkinUpdate builder.
func (suo *SkinUpdateOne) Where(ps ...predicate.Skin) *SkinUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SkinUpdateOne) Select(field string, fields ...string) *SkinUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Skin entity.
func (suo *SkinUpdateOne) Save(ctx context.Context) (*Skin, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SkinUpdateOne) SaveX(ctx context.Context) *Skin {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SkinUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SkinUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SkinUpdateOne) check() error {
	if _, ok := suo.mutation.UserID(); suo.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Skin.user"`)
	}
	return nil
}

func (suo *SkinUpdateOne) sqlSave(ctx context.Context) (_node *Skin, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(skin.Table, skin.Columns, sqlgraph.NewFieldSpec(skin.FieldID, field.TypeInt))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Skin.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, skin.FieldID)
		for _, f := range fields {
			if !skin.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != skin.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.SkinHash(); ok {
		_spec.SetField(skin.FieldSkinHash, field.TypeString, value)
	}
	if value, ok := suo.mutation.GetType(); ok {
		_spec.SetField(skin.FieldType, field.TypeUint8, value)
	}
	if value, ok := suo.mutation.AddedType(); ok {
		_spec.AddField(skin.FieldType, field.TypeUint8, value)
	}
	if value, ok := suo.mutation.Variant(); ok {
		_spec.SetField(skin.FieldVariant, field.TypeString, value)
	}
	if suo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   skin.UserTable,
			Columns: []string{skin.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   skin.UserTable,
			Columns: []string{skin.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Skin{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{skin.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}