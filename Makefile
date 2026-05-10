.PHONY: build test cover run deps compose-up compose-down mocks

build:
	go build -o bin/gophermart ./cmd/gophermart

test:
	go test -race -count=1 ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

run:
	go run ./cmd/gophermart -a localhost:8888 -d "postgres://user:password@localhost:5432/gophermart?sslmode=disable" -r "http://localhost:8081"

deps:
	go mod tidy

compose-up:
	docker compose up -d db

compose-down:
	docker compose down --remove-orphans

mocks:
	go install go.uber.org/mock/mockgen@v0.6.0
	PATH="$$PATH:$$(go env GOPATH)/bin" go generate ./internal/mocks
