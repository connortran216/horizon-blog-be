.PHONY: docs dev build

docs:
	swag init

dev: docs
	go run main.go

build: docs
	go build -o bin/go-crud main.go

test:
	go test ./...

gotestsum:
	gotestsum --format=short-verbose ./test -- -count=1 -v