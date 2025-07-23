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
	goose -dir db/migrations create init_autotm_admin sql

migrate_up:
	goose -dir db/migrations postgres "$(DB_URL)" up

migrate_down:
	goose -dir db/migrations postgres "$(DB_URL)" down

migrate_status:
	goose -dir db/migrations postgres "$(DB_URL)" status