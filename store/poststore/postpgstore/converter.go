package postpgstore

import (
	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Converts from db.Post to model.Post
func toModelPost(dbPost db.Post) model.Post {
	return model.Post{
		ID:        dbPost.ID,
		UserID:    fromPgTypeUUID(dbPost.UserID),
		Title:     dbPost.Title,
		Content:   dbPost.Content,
		CreatedAt: dbPost.CreatedAt.Time,
		UpdatedAt: dbPost.UpdatedAt.Time,
	}
}

// Converts from db.Post slice to model.Post slice
func toModelPostList(dbPosts []db.Post) []model.Post {
	posts := make([]model.Post, len(dbPosts))
	for i, dbPost := range dbPosts {
		posts[i] = toModelPost(dbPost)
	}
	return posts
}

// Converts uuid.UUID to pgtype.UUID
func toPgTypeUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: id, Valid: true}
}

// Converts pgtype.UUID to uuid.UUID
func fromPgTypeUUID(id pgtype.UUID) uuid.UUID {
	if !id.Valid {
		return uuid.Nil
	}
	return id.Bytes
}
