package model

import "github.com/google/uuid"

type Feed struct {
	ID    uuid.UUID
	URL   string
	Title string
}
