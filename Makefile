BBINARY_NAME=app
.PHONY: build clean run_app_go run generate init_swagger dev deps test update pull

build:
	go build -o ./cmd/autotm-admin ./cmd/main.go

# Clean up
clean:
	go clean
	rm -f ./cmd/autotm-admin

run_app_go:
	go run cmd/main.go

run:
	./cmd/autotm-admin

generate:
	go run ./cmd/generate/generate.go

init_swagger:
	$(GOPATH)/bin/swag init --dir ./ -g $(SRC_DIR)/cmd/main.go

dev:
	$(shell go env GOPATH)/bin/swag init --dir ./ -g $(SRC_DIR)/cmd/main.go
	go run ./cmd/generate/generate.go
	go run cmd/main.go

deps:
	go mod tidy
	#$(GOGET) github.com/example/dependency

test:
	go test -v ./...

DB_URL := postgres://autotm:autotm@127.0.0.1:5432/autotm_admin?sslmode=disable

migrate_create:
	$(shell go env GOPATH)/bin/migrate create -ext sql -dir db/migrations -seq init_autotm_admin

migrate_up:
	$(shell go env GOPATH)/bin/migrate -path db/migrations -database "$(DB_URL)" up

migrate_down:
	$(shell go env GOPATH)/bin/migrate -path db/migrations -database "$(DB_URL)" down

migrate_fix:
	@VERSION=$$(psql "$(DB_URL)" -Atc "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"); \
	echo "Automatic version: $$VERSION"; \
	$(shell go env GOPATH)/bin/migrate -path db/migrations -database "$(DB_URL)" force $$VERSION