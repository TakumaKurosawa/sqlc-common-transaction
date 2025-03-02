.PHONY: db-up db-down db-reset sqlc lint test

# Database commands
db-up:
	docker-compose up -d

db-down:
	docker-compose down

db-reset:
	docker-compose down -v
	docker-compose up -d

# SQLC commands
sqlc:
	sqlc generate

# Go commands
lint:
	go vet ./...
	go fmt ./...

test:
	go test -v ./...

# Combined commands
setup: db-up sqlc
	go mod tidy -v

# Run example
run:
	go run cmd/main.go 
