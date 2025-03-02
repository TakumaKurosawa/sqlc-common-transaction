package service

import (
	"context"
	"fmt"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/transaction"
)

// Service represents the application service layer
type Service struct {
	txManager transaction.Manager
	userStore userstore.Store
	postStore poststore.Store
}

// New creates a new service with the given transaction manager and stores
func New(txManager transaction.Manager, userStore userstore.Store, postStore poststore.Store) *Service {
	return &Service{
		txManager: txManager,
		userStore: userStore,
		postStore: postStore,
	}
}

// CreateUserWithPost creates a user and a post in a single transaction
func (s *Service) CreateUserWithPost(ctx context.Context, name, email, postTitle, postContent string) (*model.User, *model.Post, error) {
	var user model.User
	var post model.Post

	err := s.txManager.ExecTx(ctx, func(ctx context.Context) error {
		var err error

		// Create user
		user, err = s.userStore.CreateUser(ctx, name, email)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Create post
		post, err = s.postStore.CreatePost(ctx, user.ID, postTitle, postContent)
		if err != nil {
			return fmt.Errorf("failed to create post: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &user, &post, nil
}
