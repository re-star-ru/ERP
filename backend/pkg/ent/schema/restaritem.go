package schema

import (
	"backend/pkg/photo"
	"backend/pkg/restaritem"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Restaritem holds the schema definition for the Restaritem entity.
type Restaritem struct {
	ent.Schema
}

// Fields of the Restaritem.
func (Restaritem) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("onecGUID"),
		field.String("name"),
		field.String("sku"),
		field.String("itemGUID"),
		field.String("charGUID"),
		field.String("description"),
		field.String("inspector"),
		field.Strings("inspection"),

		field.JSON("photos", []photo.Photo{}),
		field.JSON("works", []restaritem.Work{}),
	}
}

// Edges of the Restaritem.
func (Restaritem) Edges() []ent.Edge {
	return nil
}
