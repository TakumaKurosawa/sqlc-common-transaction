package postpgstore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/google/uuid"
)

type pgStore struct {
	q *db.Queries
}

// New creates a new PostgreSQL implementation of poststore.Store
func New(q *db.Queries) poststore.Store {
	return &pgStore{q: q}
}

func (s *pgStore) CreatePost(ctx context.Context, userID uuid.UUID, title, content string) (model.Post, error) {
	dbParams := db.CreatePostParams{
		UserID:  toPgTypeUUID(userID),
		Title:   title,
		Content: content,
	}

	dbPost, err := s.q.CreatePost(ctx, dbParams)
	if err != nil {
		return model.Post{}, err
	}

	return toModelPost(dbPost), nil
}

func (s *pgStore) GetPost(ctx context.Context, id uuid.UUID) (model.Post, error) {
	dbPost, err := s.q.GetPost(ctx, id)
	if err != nil {
		return model.Post{}, err
	}

	return toModelPost(dbPost), nil
}

func (s *pgStore) ListPostsByUser(ctx context.Context, userID uuid.UUID) ([]model.Post, error) {
	pgUserID := toPgTypeUUID(userID)
	dbPosts, err := s.q.ListPostsByUser(ctx, pgUserID)
	if err != nil {
		return nil, err
	}

	return toModelPostList(dbPosts), nil
}

func (s *pgStore) UpdatePost(ctx context.Context, id uuid.UUID, title, content string) (model.Post, error) {
	dbParams := db.UpdatePostParams{
		ID:      id,
		Title:   title,
		Content: content,
	}

	dbPost, err := s.q.UpdatePost(ctx, dbParams)
	if err != nil {
		return model.Post{}, err
	}

	return toModelPost(dbPost), nil
}

func (s *pgStore) DeletePost(ctx context.Context, id uuid.UUID) error {
	return s.q.DeletePost(ctx, id)
}
