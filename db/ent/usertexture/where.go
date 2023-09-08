// Code generated by ent, DO NOT EDIT.

package usertexture

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/xmdhs/authlib-skin/db/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLTE(FieldID, id))
}

// UserProfileID applies equality check predicate on the "user_profile_id" field. It's identical to UserProfileIDEQ.
func UserProfileID(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldUserProfileID, v))
}

// TextureID applies equality check predicate on the "texture_id" field. It's identical to TextureIDEQ.
func TextureID(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldTextureID, v))
}

// UserProfileIDEQ applies the EQ predicate on the "user_profile_id" field.
func UserProfileIDEQ(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldUserProfileID, v))
}

// UserProfileIDNEQ applies the NEQ predicate on the "user_profile_id" field.
func UserProfileIDNEQ(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNEQ(FieldUserProfileID, v))
}

// UserProfileIDIn applies the In predicate on the "user_profile_id" field.
func UserProfileIDIn(vs ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldIn(FieldUserProfileID, vs...))
}

// UserProfileIDNotIn applies the NotIn predicate on the "user_profile_id" field.
func UserProfileIDNotIn(vs ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNotIn(FieldUserProfileID, vs...))
}

// TextureIDEQ applies the EQ predicate on the "texture_id" field.
func TextureIDEQ(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldTextureID, v))
}

// TextureIDNEQ applies the NEQ predicate on the "texture_id" field.
func TextureIDNEQ(v int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNEQ(FieldTextureID, v))
}

// TextureIDIn applies the In predicate on the "texture_id" field.
func TextureIDIn(vs ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldIn(FieldTextureID, vs...))
}

// TextureIDNotIn applies the NotIn predicate on the "texture_id" field.
func TextureIDNotIn(vs ...int) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNotIn(FieldTextureID, vs...))
}

// HasUserProfile applies the HasEdge predicate on the "user_profile" edge.
func HasUserProfile() predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, UserProfileTable, UserProfileColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserProfileWith applies the HasEdge predicate on the "user_profile" edge with a given conditions (other predicates).
func HasUserProfileWith(preds ...predicate.UserProfile) predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		step := newUserProfileStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTexture applies the HasEdge predicate on the "texture" edge.
func HasTexture() predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, TextureTable, TextureColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTextureWith applies the HasEdge predicate on the "texture" edge with a given conditions (other predicates).
func HasTextureWith(preds ...predicate.Texture) predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		step := newTextureStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.UserTexture) predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.UserTexture) predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.UserTexture) predicate.UserTexture {
	return predicate.UserTexture(func(s *sql.Selector) {
		p(s.Not())
	})
}
