// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID        uuid.UUID        `json:"id"`
	UserID    pgtype.UUID      `json:"userId"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}

type User struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`
}
