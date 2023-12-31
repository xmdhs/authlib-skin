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

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldType, v))
}

// Variant applies equality check predicate on the "variant" field. It's identical to VariantEQ.
func Variant(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldVariant, v))
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

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldContainsFold(FieldType, v))
}

// VariantEQ applies the EQ predicate on the "variant" field.
func VariantEQ(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEQ(FieldVariant, v))
}

// VariantNEQ applies the NEQ predicate on the "variant" field.
func VariantNEQ(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNEQ(FieldVariant, v))
}

// VariantIn applies the In predicate on the "variant" field.
func VariantIn(vs ...string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldIn(FieldVariant, vs...))
}

// VariantNotIn applies the NotIn predicate on the "variant" field.
func VariantNotIn(vs ...string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldNotIn(FieldVariant, vs...))
}

// VariantGT applies the GT predicate on the "variant" field.
func VariantGT(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGT(FieldVariant, v))
}

// VariantGTE applies the GTE predicate on the "variant" field.
func VariantGTE(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldGTE(FieldVariant, v))
}

// VariantLT applies the LT predicate on the "variant" field.
func VariantLT(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLT(FieldVariant, v))
}

// VariantLTE applies the LTE predicate on the "variant" field.
func VariantLTE(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldLTE(FieldVariant, v))
}

// VariantContains applies the Contains predicate on the "variant" field.
func VariantContains(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldContains(FieldVariant, v))
}

// VariantHasPrefix applies the HasPrefix predicate on the "variant" field.
func VariantHasPrefix(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldHasPrefix(FieldVariant, v))
}

// VariantHasSuffix applies the HasSuffix predicate on the "variant" field.
func VariantHasSuffix(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldHasSuffix(FieldVariant, v))
}

// VariantEqualFold applies the EqualFold predicate on the "variant" field.
func VariantEqualFold(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldEqualFold(FieldVariant, v))
}

// VariantContainsFold applies the ContainsFold predicate on the "variant" field.
func VariantContainsFold(v string) predicate.UserTexture {
	return predicate.UserTexture(sql.FieldContainsFold(FieldVariant, v))
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
