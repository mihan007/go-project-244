test:
	go test -v ./...

deps:
	go mod tidy

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run ./...

coverage:
	go test -coverprofile=coverage.out