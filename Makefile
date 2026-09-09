.PHONY: install start test

install:
	go mod tidy
	docker compose up -d

start:
	go run cmd/main.go

test:
	go test ./... -v

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html