include help.mk

APP_NAME?=filsorter
BUILD_DIR=./

clean: # Clean generated files and test cache
	@rm -f ${BUILD_DIR}${APP_NAME}
	@go clean -testcache

fmt: # Format Go source code
	@go fmt ./...

test: clean # Run unit tests
	@go test -v -count=1 -cover ./...

race: # Run unit tests with race
	go test -v -race -count=1 ./...

.PHONY: cover 
cover: clean # Show test coverage
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: lint 
lint: clean # Run linters
	golangci-lint run

.PHONY: build
build: clean # Build binary
	@env GOOS=linux go build -o ${BUILD_DIR}${APP_NAME} .
