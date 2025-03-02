# Common transaction interface sample

A Go package for common transaction handling using SQLC.

## Overview

This package provides a common interface for handling database transactions in Go applications.
It supports both the standard `database/sql` and PostgreSQL's `pgx` driver.
The main goal is to abstract and standardize transaction handling across different database implementations,
which significantly improves the testability of your application code.

## Features

- Automatic management of transaction boundaries
- Implementations for `database/sql` and `pgx`
- Context-based transaction sharing
- Automatic rollback on error
- Interface-based design for easy mocking in tests
- Improved testability for business logic without database dependencies

## Development

### Setting Up Local PostgreSQL with Docker

This repository includes a Docker Compose configuration for easy local development with PostgreSQL.

1. Start the PostgreSQL container:

   ```bash
   docker-compose up -d
   # or
   make db-up
   ```

2. The database will be initialized with the schema defined in `sql/schema.sql` and sample data from `sql/seed.sql`.

3. To stop the container:

   ```bash
   docker-compose down
   # or
   make db-down
   ```

4. To stop the container and remove volumes (data will be lost):

   ```bash
   docker-compose down -v
   # or
   make db-reset
   ```

Connection details:

- Host: localhost
- Port: 5432
- User: postgres
- Password: postgres
- Database: postgres

### Using Makefile

This repository includes a Makefile with useful commands for development:

```bash
# Start the database
make db-up

# Stop the database
make db-down

# Reset the database (down + up)
make db-reset

# Generate SQLC code
make sqlc

# Run linting
make lint

# Run tests
make test

# Complete setup (db-up + sqlc + go mod tidy)
make setup

# Run the example application
make run
```

## Transaction Abstraction Design

This package demonstrates how to abstract database transaction handling, hiding the specific implementations like `database/sql` or `pgx` behind a common interface. Here's how the abstraction works:

### Interface-Based Design

The core of the abstraction is the `Manager` interface in the `transaction` package:

```go
type Manager interface {
    Begin(ctx context.Context) (context.Context, error)
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
    ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

This interface defines the contract for all transaction managers, regardless of the underlying database driver.

### Context-Based Transaction Sharing

Transactions are stored in and retrieved from the context:

1. When a transaction begins, it's stored in the context using a driver-specific key
2. The context with the transaction is passed to all functions that need to use the transaction
3. Functions retrieve the transaction from the context when needed

This approach prevents having to pass transaction objects through multiple layers of function calls.

### Implementation Hiding

The package includes two implementations of the `Manager` interface:

1. **SQL Implementation** (`sqltransaction.Manager`):

   - Wraps `database/sql.DB` and `sql.Tx`
   - Provides methods to begin, commit, and rollback transactions
   - Uses context to store and retrieve `sql.Tx` objects

2. **PGX Implementation** (`pgxtransaction.Manager`):
   - Wraps `pgx/v5/pgxpool.Pool` and `pgx.Tx`
   - Provides the same interface but uses PGX's transaction types
   - Uses context to store and retrieve `pgx.Tx` objects

Each implementation handles its specific driver details internally, while exposing the same interface to callers.

### Transaction Retrieval

To use a transaction, services retrieve it from the context:

```go
// For SQL transactions
tx, err := sqltransaction.GetTx(ctx)

// For PGX transactions
tx, err := pgxtransaction.GetTx(ctx)
```

Services can then use these transactions with their database operations, but they don't need to know how the transaction was created or how it will be committed/rolled back.

### Benefits of this Abstraction

1. **Separation of Concerns**: Transaction management is separate from business logic
2. **Testability**: Easy to mock the transaction interface for testing
3. **Driver Independence**: Services don't need to know which database driver is being used
4. **Consistency**: Standardized transaction handling across the application
5. **Error Handling**: Automatic rollback on error provides safety
6. **Future Flexibility**: New database drivers can be supported by adding new implementations

## Testing with Mock

One of the key benefits of the abstract transaction interface is improved testability. By using the interface-based design,
you can easily mock the transaction manager and focus on testing your business logic without database dependencies.

### Creating Mocks

This project uses [uber-go/mock](https://github.com/uber-go/mock) for generating mocks.

1. Install mockgen:

   ```bash
   go install go.uber.org/mock/mockgen@latest
   ```

2. Add go:generate comments to your interface files:

   ```go
   // In transaction/transaction.go
   //go:generate mockgen -destination=mocks/mock_transaction.go -package=mocks github.com/TakumaKurosawa/sqlc-common-transaction/transaction Manager

   // In store/userstore/store.go
   //go:generate mockgen -destination=mocks/mock_store.go -package=mocks github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore Store

   // In store/poststore/store.go
   //go:generate mockgen -destination=mocks/mock_store.go -package=mocks github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore Store
   ```

3. Generate all mocks at once:

   ```bash
   go generate ./...
   ```

### Writing Tests with Mocks

Here's an example of a table-driven test using mocks:

```go
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
        // More test cases...
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
```

### Benefits of This Testing Approach

1. **Focus on Business Logic**: By mocking the database interactions, tests can focus on the business logic without real database connections.
2. **Faster Tests**: No need to set up database fixtures or wait for database operations.
3. **Isolated Tests**: Each test runs in isolation, reducing test flakiness.
4. **Predictable Results**: Tests provide consistent results regardless of database state.
5. **Comprehensive Coverage**: Easy to test various scenarios including error cases.

## Installation

```bash
git clone github.com/TakumaKurosawa/sqlc-common-transaction
go mod tidy -v  # Update dependencies
```

## Requirements

- Go 1.18 or higher
- PostgreSQL 10 or higher (when using pgx driver)
- Docker and Docker Compose (for local development)

## License

This project is released under the MIT License.
