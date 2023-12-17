# packages excluded from tests: config
TEST_PACKAGES=$(shell go list ./... | grep -v "config")
COVERAGE_FILE=coverage.out
SERVICE_NAME=whatsapp-openai

run:
	@export $(shell cat .env | xargs) && go run -race .

test:
	go test -race $(TEST_PACKAGES) -count=1

format:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run ./...

vet:
	go vet ./...

security:
	gosec ./...

coverage:
	go test $(TEST_PACKAGES) -coverprofile=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE)

build: format lint vet test
	go mod tidy
	go build -o bin/$(SERVICE_NAME)
	@echo "Build completed"

clean:
	rm -rf ./bin
	rm -f coverage.out

ci:
	go test -race $(TEST_PACKAGES) -coverprofile=$(COVERAGE_FILE)

.PHONY: run test format lint vet security coverage build clean ci