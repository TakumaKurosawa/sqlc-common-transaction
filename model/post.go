package model

import (
	"time"

	"github.com/google/uuid"
)

// Post represents a post entity in the domain
type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
