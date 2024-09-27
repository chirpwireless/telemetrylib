all: test build

install:
	@go mod download

build: install
	@go build

test: build
	@go test -v ./...
