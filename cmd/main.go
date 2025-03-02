package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TakumaKurosawa/sqlc-common-transaction/pkg/db"
	"github.com/TakumaKurosawa/sqlc-common-transaction/service"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/poststore/postpgstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/store/userstore/userpgstore"
	"github.com/TakumaKurosawa/sqlc-common-transaction/transaction/pgxtransaction"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/example")
	if err != nil {
		log.Fatalf("Unable to parse connection URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	txManager := pgxtransaction.New(pool)
	queries := db.New(pool)

	userStore := userpgstore.New(queries)
	postStore := postpgstore.New(queries)

	svc := service.New(txManager, userStore, postStore)

	user, post, err := svc.CreateUserWithPost(context.Background(), "John Doe", "john@example.com", "First Post", "Hello, World!")
	if err != nil {
		log.Fatalf("Failed to create user with post: %v", err)
	}

	fmt.Printf("Created user: %s (%s)\n", user.Name, user.ID.String())
	fmt.Printf("Created post: %s (%s)\n", post.Title, post.ID.String())
}
