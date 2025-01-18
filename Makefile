test:
	go test
mod:
	go mod download
build: mod
	go build -o bin/link-checker main.go