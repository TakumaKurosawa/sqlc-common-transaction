package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/transaction"
	"github.com/jackc/pgx/v5/pgtype"
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

	// Create a query object to interact with the database
	queries := db.New(pool)

	txManager := transaction.NewPgxManager(pool)

	// Example of processing within a transaction
	if err := txManager.ExecTx(context.Background(), func(ctx context.Context) error {
		// Example: Create a user and create a post associated with that user

		tx, err := transaction.GetPgxTx(ctx)
		if err != nil {
			return err
		}

		// Use queries with transaction context
		q := queries.WithTx(tx)

		// Create user using SQLC generated function
		user, err := q.CreateUser(ctx, db.CreateUserParams{
			Name:  "Sample User",
			Email: "sample@example.com",
		})
		if err != nil {
			return fmt.Errorf("user creation error: %w", err)
		}

		// Convert uuid.UUID to pgtype.UUID
		pgUserID := pgtype.UUID{}
		pgUserID.Bytes = user.ID
		pgUserID.Valid = true

		// Create post using SQLC generated function
		if _, err := q.CreatePost(ctx, db.CreatePostParams{
			UserID:  pgUserID,
			Title:   "Sample Post",
			Content: "This is a sample post content.",
		}); err != nil {
			return fmt.Errorf("post creation error: %w", err)
		}

		fmt.Printf("Created user with ID %s and added a related post\n", user.ID)
		return nil
	}); err != nil {
		log.Fatalf("Transaction execution error: %v", err)
	}

	fmt.Println("Process completed successfully")
}
