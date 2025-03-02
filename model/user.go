package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity in the domain
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
