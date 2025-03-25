package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Article struct {
	ent.Schema
}

func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").Positive(),
		field.String("title"),
		field.Text("content"),
		field.String("url").Unique(),
		field.String("author"),
		field.String("source"),
		field.String("summary").Optional(),
		field.JSON("tags", []string{}).Optional(),
		field.Time("published_at"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Article) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
		index.Fields("url"),
		index.Fields("author"),
		index.Fields("published_at"),
	}
}

func (Article) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("articles").
			Unique().
			Required(),
	}
}
