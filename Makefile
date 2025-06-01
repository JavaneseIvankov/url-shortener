include .env
export 

.PHONY: air migrate migrate-new

DATABASE_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=disable

run:
	go run ./cmd/app/

hr:
	air

lint:
	golangci-lint run

migrate:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

migrate.new:
	@read -p "Enter name (e.g. add_index): " name; \
	migrate create -dir db/migrations -ext sql $${name} \ 
	
migrate.%:
	migrate -path db/migrations -database "$(DATABASE_URL)" $*
