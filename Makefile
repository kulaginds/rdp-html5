test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/server cmd/server/main.go
