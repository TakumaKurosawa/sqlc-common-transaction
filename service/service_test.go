package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TakumaKurosawa/sqlc-common-transaction/model"
	postmocks "github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore/mocks"
	usermocks "github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore/mocks"
	txmocks "github.com/TakumaKurosawa/sqlc-common-transaction/transaction/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUserWithPost(t *testing.T) {
	userID := uuid.New()
	postID := uuid.New()
	now := time.Now()

	testUser := model.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	testPost := model.Post{
		ID:        postID,
		UserID:    userID,
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := map[string]struct {
		setupMocks      func(mockTx *txmocks.MockManager, mockUserStore *usermocks.MockStore, mockPostStore *postmocks.MockStore)
		userName        string
		userEmail       string
		postTitle       string
		postContent     string
		expectedUser    *model.User
		expectedPost    *model.Post
		expectedError   assert.ErrorAssertionFunc
		expectedErrText string
	}{
		"success - user and post created successfully": {
			setupMocks: func(mockTx *txmocks.MockManager, mockUserStore *usermocks.MockStore, mockPostStore *postmocks.MockStore) {
				mockTx.EXPECT().
					ExecTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					})

				mockUserStore.EXPECT().
					CreateUser(gomock.Any(), "Test User", "test@example.com").
					Return(testUser, nil)

				mockPostStore.EXPECT().
					CreatePost(gomock.Any(), userID, "Test Title", "Test Content").
					Return(testPost, nil)
			},
			userName:        "Test User",
			userEmail:       "test@example.com",
			postTitle:       "Test Title",
			postContent:     "Test Content",
			expectedUser:    &testUser,
			expectedPost:    &testPost,
			expectedError:   assert.NoError,
			expectedErrText: "",
		},
		"error - user creation fails": {
			setupMocks: func(mockTx *txmocks.MockManager, mockUserStore *usermocks.MockStore, mockPostStore *postmocks.MockStore) {
				mockTx.EXPECT().
					ExecTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).
					Return(errors.New("failed to create user"))

				mockUserStore.EXPECT().
					CreateUser(gomock.Any(), "Test User", "test@example.com").
					Return(model.User{}, errors.New("user creation failed"))
			},
			userName:        "Test User",
			userEmail:       "test@example.com",
			postTitle:       "Test Title",
			postContent:     "Test Content",
			expectedUser:    nil,
			expectedPost:    nil,
			expectedError:   assert.Error,
			expectedErrText: "failed to create user",
		},
		"error - post creation fails": {
			setupMocks: func(mockTx *txmocks.MockManager, mockUserStore *usermocks.MockStore, mockPostStore *postmocks.MockStore) {
				mockTx.EXPECT().
					ExecTx(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
						return fn(ctx)
					}).
					Return(errors.New("failed to create post"))

				mockUserStore.EXPECT().
					CreateUser(gomock.Any(), "Test User", "test@example.com").
					Return(testUser, nil)

				mockPostStore.EXPECT().
					CreatePost(gomock.Any(), userID, "Test Title", "Test Content").
					Return(model.Post{}, errors.New("post creation failed"))
			},
			userName:        "Test User",
			userEmail:       "test@example.com",
			postTitle:       "Test Title",
			postContent:     "Test Content",
			expectedUser:    nil,
			expectedPost:    nil,
			expectedError:   assert.Error,
			expectedErrText: "failed to create post",
		},
		"error - transaction execution fails": {
			setupMocks: func(mockTx *txmocks.MockManager, mockUserStore *usermocks.MockStore, mockPostStore *postmocks.MockStore) {
				mockTx.EXPECT().
					ExecTx(gomock.Any(), gomock.Any()).
					Return(errors.New("transaction error"))
			},
			userName:        "Test User",
			userEmail:       "test@example.com",
			postTitle:       "Test Title",
			postContent:     "Test Content",
			expectedUser:    nil,
			expectedPost:    nil,
			expectedError:   assert.Error,
			expectedErrText: "transaction error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTx := txmocks.NewMockManager(ctrl)
			mockUserStore := usermocks.NewMockStore(ctrl)
			mockPostStore := postmocks.NewMockStore(ctrl)

			tt.setupMocks(mockTx, mockUserStore, mockPostStore)

			svc := New(mockTx, mockUserStore, mockPostStore)

			user, post, err := svc.CreateUserWithPost(context.Background(), tt.userName, tt.userEmail, tt.postTitle, tt.postContent)

			tt.expectedError(t, err)
			if tt.expectedErrText != "" {
				assert.Contains(t, err.Error(), tt.expectedErrText)
			}
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedPost, post)
		})
	}
}
