# sqlc-common-transaction

A Go package for common transaction handling using SQLC.

## Overview

This package provides a common interface for handling database transactions in Go applications.
It supports both the standard `database/sql` and PostgreSQL's `pgx` driver.

## Features

- Automatic management of transaction boundaries
- Implementations for `database/sql` and `pgx`
- Context-based transaction sharing
- Automatic rollback on error

## Usage

### Initializing Transaction Manager

```go
// For database/sql
db, _ := sql.Open("postgres", "postgres://user:pass@localhost/db")
txManager := transaction.NewTxManager(db)

// For pgx
pool, _ := pgxpool.New(context.Background(), "postgres://user:pass@localhost/db")
pgxManager := transaction.NewPgxManager(pool)
```

### Executing Transactions

```go
err := txManager.ExecTx(ctx, func(ctx context.Context) error {
    // Write the process to be executed within the transaction here
    // If an error is returned, it will be rolled back automatically
    return nil
})
```

### Getting Transaction

```go
// For database/sql
tx, err := transaction.GetTx(ctx)
if err != nil {
    return err
}

// For pgx
pgxTx, err := transaction.GetPgxTx(ctx)
if err != nil {
    return err
}
```

## Combining with SQLC

This package is optimal for use with code generated by [SQLC](https://github.com/kyleconroy/sqlc).
You can pass transactions to query functions generated by SQLC.

### Using SQLC with This Package

1. Install SQLC:

   ```bash
   go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
   ```

2. Generate code using SQLC:

   ```bash
   sqlc generate
   ```

3. Update dependencies:

   ```bash
   go mod tidy -v
   ```

4. Use the generated code with transactions:

   ```go
   // Create a query object
   queries := db.New(pool)

   // Execute in transaction
   if err := txManager.ExecTx(ctx, func(ctx context.Context) error {
       tx, err := transaction.GetPgxTx(ctx)
       if err != nil {
           return err
       }

       // Use queries with transaction
       q := queries.WithTx(tx)

       // Call generated methods
       user, err := q.CreateUser(ctx, params)
       if err != nil {
           return err
       }

       // More operations...
       return nil
   }); err != nil {
       // Handle error
   }
   ```

## Installation

```bash
go get github.com/TakumaKurosawa/sqlc-common-transaction
go mod tidy -v  # Update dependencies
```

## Requirements

- Go 1.18 or higher
- PostgreSQL 10 or higher (when using pgx driver)

## License

This project is released under the MIT License.
