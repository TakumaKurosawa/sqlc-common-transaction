package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/transaction"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	txManager := transaction.NewPgxManager(pool)

	// Example of processing within a transaction
	if err := txManager.ExecTx(context.Background(), func(ctx context.Context) error {
		// Example: Create a user and create a post associated with that user

		tx, err := transaction.GetPgxTx(ctx)
		if err != nil {
			return err
		}

		var userID uuid.UUID
		if err := tx.QueryRow(ctx,
			"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
			"Sample User", "sample@example.com").Scan(&userID); err != nil {
			return fmt.Errorf("user creation error: %w", err)
		}

		if _, err := tx.Exec(ctx,
			"INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)",
			userID, "Sample Post", "This is a sample post content."); err != nil {
			return fmt.Errorf("post creation error: %w", err)
		}

		fmt.Printf("Created user with ID %s and added a related post\n", userID)
		return nil
	}); err != nil {
		log.Fatalf("Transaction execution error: %v", err)
	}

	fmt.Println("Process completed successfully")
}
