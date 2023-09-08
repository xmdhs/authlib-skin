// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/xmdhs/authlib-skin/db/ent/predicate"
	"github.com/xmdhs/authlib-skin/db/ent/texture"
)

// TextureDelete is the builder for deleting a Texture entity.
type TextureDelete struct {
	config
	hooks    []Hook
	mutation *TextureMutation
}

// Where appends a list predicates to the TextureDelete builder.
func (td *TextureDelete) Where(ps ...predicate.Texture) *TextureDelete {
	td.mutation.Where(ps...)
	return td
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (td *TextureDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, td.sqlExec, td.mutation, td.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (td *TextureDelete) ExecX(ctx context.Context) int {
	n, err := td.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (td *TextureDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(texture.Table, sqlgraph.NewFieldSpec(texture.FieldID, field.TypeInt))
	if ps := td.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, td.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	td.mutation.done = true
	return affected, err
}

// TextureDeleteOne is the builder for deleting a single Texture entity.
type TextureDeleteOne struct {
	td *TextureDelete
}

// Where appends a list predicates to the TextureDelete builder.
func (tdo *TextureDeleteOne) Where(ps ...predicate.Texture) *TextureDeleteOne {
	tdo.td.mutation.Where(ps...)
	return tdo
}

// Exec executes the deletion query.
func (tdo *TextureDeleteOne) Exec(ctx context.Context) error {
	n, err := tdo.td.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{texture.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tdo *TextureDeleteOne) ExecX(ctx context.Context) {
	if err := tdo.Exec(ctx); err != nil {
		panic(err)
	}
}
