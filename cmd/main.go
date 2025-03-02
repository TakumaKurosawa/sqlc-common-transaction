package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TakumaKurosawa/sqlc-common-transaction/store"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Parse database config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Create connection pool: %v", err)
	}
	defer pool.Close()

	// Create transaction manager using PostgreSQL implementation
	txManager := store.NewTxManager(pool)

	// Execute transaction across multiple tables
	if err := txManager.ExecTx(context.Background(), func(userStore userstore.Store, postStore poststore.Store) error {
		user, err := userStore.CreateUser(context.Background(), "Alice", "alice@example.com")
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		fmt.Printf("Created user: %v\n", user)

		post, err := postStore.CreatePost(context.Background(), user.ID, "My First Post", "Hello, World!")
		if err != nil {
			return fmt.Errorf("create post: %w", err)
		}
		fmt.Printf("Created post: %v\n", post)

		return nil
	}); err != nil {
		log.Fatalf("Transaction error: %v", err)
	}
}
