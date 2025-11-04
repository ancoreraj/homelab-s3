.PHONY: run build clean test deps dev

run:
	go run cmd/server/main.go

build:
	go build -o bin/s3-clone cmd/server/main.go

clean:
	rm -rf bin/
	rm -rf tmp/

deps:
	go mod download
	go mod tidy

dev:
	air -c .air.toml
fmt:
	go fmt ./...
lint:
	golangci-lint run
test:
	go test -v ./...

build-all:
	GOOS=darwin GOARCH=amd64 go build -o bin/s3-clone-darwin-amd64 cmd/server/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/s3-clone-darwin-arm64 cmd/server/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/s3-clone-linux-amd64 cmd/server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/s3-clone-windows-amd64.exe cmd/server/main.go