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
	"github.com/xmdhs/authlib-skin/db/ent/texture"
	"github.com/xmdhs/authlib-skin/db/ent/user"
	"github.com/xmdhs/authlib-skin/db/ent/userprofile"
)

// UserProfileUpdate is the builder for updating UserProfile entities.
type UserProfileUpdate struct {
	config
	hooks    []Hook
	mutation *UserProfileMutation
}

// Where appends a list predicates to the UserProfileUpdate builder.
func (upu *UserProfileUpdate) Where(ps ...predicate.UserProfile) *UserProfileUpdate {
	upu.mutation.Where(ps...)
	return upu
}

// SetName sets the "name" field.
func (upu *UserProfileUpdate) SetName(s string) *UserProfileUpdate {
	upu.mutation.SetName(s)
	return upu
}

// SetUUID sets the "uuid" field.
func (upu *UserProfileUpdate) SetUUID(s string) *UserProfileUpdate {
	upu.mutation.SetUUID(s)
	return upu
}

// SetUserID sets the "user" edge to the User entity by ID.
func (upu *UserProfileUpdate) SetUserID(id int) *UserProfileUpdate {
	upu.mutation.SetUserID(id)
	return upu
}

// SetUser sets the "user" edge to the User entity.
func (upu *UserProfileUpdate) SetUser(u *User) *UserProfileUpdate {
	return upu.SetUserID(u.ID)
}

// AddTextureIDs adds the "texture" edge to the Texture entity by IDs.
func (upu *UserProfileUpdate) AddTextureIDs(ids ...int) *UserProfileUpdate {
	upu.mutation.AddTextureIDs(ids...)
	return upu
}

// AddTexture adds the "texture" edges to the Texture entity.
func (upu *UserProfileUpdate) AddTexture(t ...*Texture) *UserProfileUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return upu.AddTextureIDs(ids...)
}

// Mutation returns the UserProfileMutation object of the builder.
func (upu *UserProfileUpdate) Mutation() *UserProfileMutation {
	return upu.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (upu *UserProfileUpdate) ClearUser() *UserProfileUpdate {
	upu.mutation.ClearUser()
	return upu
}

// ClearTexture clears all "texture" edges to the Texture entity.
func (upu *UserProfileUpdate) ClearTexture() *UserProfileUpdate {
	upu.mutation.ClearTexture()
	return upu
}

// RemoveTextureIDs removes the "texture" edge to Texture entities by IDs.
func (upu *UserProfileUpdate) RemoveTextureIDs(ids ...int) *UserProfileUpdate {
	upu.mutation.RemoveTextureIDs(ids...)
	return upu
}

// RemoveTexture removes "texture" edges to Texture entities.
func (upu *UserProfileUpdate) RemoveTexture(t ...*Texture) *UserProfileUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return upu.RemoveTextureIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (upu *UserProfileUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, upu.sqlSave, upu.mutation, upu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (upu *UserProfileUpdate) SaveX(ctx context.Context) int {
	affected, err := upu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (upu *UserProfileUpdate) Exec(ctx context.Context) error {
	_, err := upu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (upu *UserProfileUpdate) ExecX(ctx context.Context) {
	if err := upu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (upu *UserProfileUpdate) check() error {
	if _, ok := upu.mutation.UserID(); upu.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserProfile.user"`)
	}
	return nil
}

func (upu *UserProfileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := upu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(userprofile.Table, userprofile.Columns, sqlgraph.NewFieldSpec(userprofile.FieldID, field.TypeInt))
	if ps := upu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := upu.mutation.Name(); ok {
		_spec.SetField(userprofile.FieldName, field.TypeString, value)
	}
	if value, ok := upu.mutation.UUID(); ok {
		_spec.SetField(userprofile.FieldUUID, field.TypeString, value)
	}
	if upu.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   userprofile.UserTable,
			Columns: []string{userprofile.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upu.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   userprofile.UserTable,
			Columns: []string{userprofile.UserColumn},
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
	if upu.mutation.TextureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upu.mutation.RemovedTextureIDs(); len(nodes) > 0 && !upu.mutation.TextureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upu.mutation.TextureIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, upu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{userprofile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	upu.mutation.done = true
	return n, nil
}

// UserProfileUpdateOne is the builder for updating a single UserProfile entity.
type UserProfileUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserProfileMutation
}

// SetName sets the "name" field.
func (upuo *UserProfileUpdateOne) SetName(s string) *UserProfileUpdateOne {
	upuo.mutation.SetName(s)
	return upuo
}

// SetUUID sets the "uuid" field.
func (upuo *UserProfileUpdateOne) SetUUID(s string) *UserProfileUpdateOne {
	upuo.mutation.SetUUID(s)
	return upuo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (upuo *UserProfileUpdateOne) SetUserID(id int) *UserProfileUpdateOne {
	upuo.mutation.SetUserID(id)
	return upuo
}

// SetUser sets the "user" edge to the User entity.
func (upuo *UserProfileUpdateOne) SetUser(u *User) *UserProfileUpdateOne {
	return upuo.SetUserID(u.ID)
}

// AddTextureIDs adds the "texture" edge to the Texture entity by IDs.
func (upuo *UserProfileUpdateOne) AddTextureIDs(ids ...int) *UserProfileUpdateOne {
	upuo.mutation.AddTextureIDs(ids...)
	return upuo
}

// AddTexture adds the "texture" edges to the Texture entity.
func (upuo *UserProfileUpdateOne) AddTexture(t ...*Texture) *UserProfileUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return upuo.AddTextureIDs(ids...)
}

// Mutation returns the UserProfileMutation object of the builder.
func (upuo *UserProfileUpdateOne) Mutation() *UserProfileMutation {
	return upuo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (upuo *UserProfileUpdateOne) ClearUser() *UserProfileUpdateOne {
	upuo.mutation.ClearUser()
	return upuo
}

// ClearTexture clears all "texture" edges to the Texture entity.
func (upuo *UserProfileUpdateOne) ClearTexture() *UserProfileUpdateOne {
	upuo.mutation.ClearTexture()
	return upuo
}

// RemoveTextureIDs removes the "texture" edge to Texture entities by IDs.
func (upuo *UserProfileUpdateOne) RemoveTextureIDs(ids ...int) *UserProfileUpdateOne {
	upuo.mutation.RemoveTextureIDs(ids...)
	return upuo
}

// RemoveTexture removes "texture" edges to Texture entities.
func (upuo *UserProfileUpdateOne) RemoveTexture(t ...*Texture) *UserProfileUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return upuo.RemoveTextureIDs(ids...)
}

// Where appends a list predicates to the UserProfileUpdate builder.
func (upuo *UserProfileUpdateOne) Where(ps ...predicate.UserProfile) *UserProfileUpdateOne {
	upuo.mutation.Where(ps...)
	return upuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (upuo *UserProfileUpdateOne) Select(field string, fields ...string) *UserProfileUpdateOne {
	upuo.fields = append([]string{field}, fields...)
	return upuo
}

// Save executes the query and returns the updated UserProfile entity.
func (upuo *UserProfileUpdateOne) Save(ctx context.Context) (*UserProfile, error) {
	return withHooks(ctx, upuo.sqlSave, upuo.mutation, upuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (upuo *UserProfileUpdateOne) SaveX(ctx context.Context) *UserProfile {
	node, err := upuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (upuo *UserProfileUpdateOne) Exec(ctx context.Context) error {
	_, err := upuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (upuo *UserProfileUpdateOne) ExecX(ctx context.Context) {
	if err := upuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (upuo *UserProfileUpdateOne) check() error {
	if _, ok := upuo.mutation.UserID(); upuo.mutation.UserCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "UserProfile.user"`)
	}
	return nil
}

func (upuo *UserProfileUpdateOne) sqlSave(ctx context.Context) (_node *UserProfile, err error) {
	if err := upuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(userprofile.Table, userprofile.Columns, sqlgraph.NewFieldSpec(userprofile.FieldID, field.TypeInt))
	id, ok := upuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "UserProfile.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := upuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, userprofile.FieldID)
		for _, f := range fields {
			if !userprofile.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != userprofile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := upuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := upuo.mutation.Name(); ok {
		_spec.SetField(userprofile.FieldName, field.TypeString, value)
	}
	if value, ok := upuo.mutation.UUID(); ok {
		_spec.SetField(userprofile.FieldUUID, field.TypeString, value)
	}
	if upuo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   userprofile.UserTable,
			Columns: []string{userprofile.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upuo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   userprofile.UserTable,
			Columns: []string{userprofile.UserColumn},
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
	if upuo.mutation.TextureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upuo.mutation.RemovedTextureIDs(); len(nodes) > 0 && !upuo.mutation.TextureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := upuo.mutation.TextureIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   userprofile.TextureTable,
			Columns: []string{userprofile.TextureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &UserProfile{config: upuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, upuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{userprofile.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	upuo.mutation.done = true
	return _node, nil
}
