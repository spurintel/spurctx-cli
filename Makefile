.PHONY: fmt lint bin all crosscompile

all: fmt lint test bin

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test ./...

bin:
	mkdir -p target
	go build -o target/spurctx


crosscompile:
	mkdir -p target
	GOOS=linux GOARCH=amd64 go build -o target/spurctx-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o target/spurctx-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o target/spurctx-windows-amd64.exe
	GOOS=linux GOARCH=arm64 go build -o target/spurctx-linux-arm64
	GOOS=darwin GOARCH=arm64 go build -o target/spurctx-darwin-arm64
