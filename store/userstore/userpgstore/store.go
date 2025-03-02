package userpgstore

import (
	"context"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/google/uuid"
)

type pgStore struct {
	q *db.Queries
}

// New creates a new PostgreSQL implementation of userstore.Store
func New(q *db.Queries) userstore.Store {
	return &pgStore{q: q}
}

func (s *pgStore) CreateUser(ctx context.Context, name, email string) (model.User, error) {
	dbParams := db.CreateUserParams{
		Name:  name,
		Email: email,
	}

	dbUser, err := s.q.CreateUser(ctx, dbParams)
	if err != nil {
		return model.User{}, err
	}

	return toModelUser(dbUser), nil
}

func (s *pgStore) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	dbUser, err := s.q.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return toModelUser(dbUser), nil
}

func (s *pgStore) ListUsers(ctx context.Context) ([]model.User, error) {
	dbUsers, err := s.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return toModelUserList(dbUsers), nil
}

func (s *pgStore) UpdateUser(ctx context.Context, id uuid.UUID, name, email string) (model.User, error) {
	dbParams := db.UpdateUserParams{
		ID:    id,
		Name:  name,
		Email: email,
	}

	dbUser, err := s.q.UpdateUser(ctx, dbParams)
	if err != nil {
		return model.User{}, err
	}

	return toModelUser(dbUser), nil
}

func (s *pgStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.q.DeleteUser(ctx, id)
}
