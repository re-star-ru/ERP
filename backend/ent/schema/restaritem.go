package schema

import (
	"backend/pkg/photo"
	"backend/pkg/work"
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
		field.String("onecGUID"),
		field.String("name").Optional(),
		field.String("sku").Optional(),
		field.String("itemGUID").Optional(),
		field.String("charGUID").Optional(),
		field.String("description").Optional(),
		field.String("inspector").Optional(),
		field.Strings("inspection").Optional(),

		field.JSON("photos", []photo.Photo{}).Optional(),
		field.JSON("works", []work.Work{}).Optional(),
	}
}

// Edges of the Restaritem.
func (Restaritem) Edges() []ent.Edge {
	return nil
}
