package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Feed struct {
	ent.Schema
}

func (Feed) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("url").NotEmpty().Unique(),
		field.String("title").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
