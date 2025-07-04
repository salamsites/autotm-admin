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
