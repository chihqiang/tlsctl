version := $(shell git describe --tags --always)
OUTPUT := tlsctl
MAIN := ./cmd/tlsctl/main.go

.PHONY: build check test clean

build:
	@echo "🔧 Building $(OUTPUT) with version $(version)..."
	GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(version)" -o $(OUTPUT) $(MAIN)
	@echo "✅ Build complete: $(OUTPUT)"


build_liunx_amd64:
	@echo "🔧 Building $(OUTPUT) with version $(version)..."
	GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=$(version)" -o $(OUTPUT) $(MAIN)
	@echo "✅ Build complete: $(OUTPUT)"

check:
	@find . -name "*.go" -exec go fmt {} \;
	@go mod tidy
	@echo "🔍 Running linters..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	golangci-lint run ./...
	@echo "✅ Linting passed"


test: build
	@echo "🧪 Running tests..."
	chmod +x tlsctl
	./tlsctl localhost
	go test -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.txt ./... -run . -timeout=2m
	@echo "🔍 Checking git status..."
	@git diff --quiet || (echo "❌ Uncommitted changes detected in working directory!" && git status && exit 1)
	@git diff --cached --quiet || (echo "❌ Staged but uncommitted changes detected!" && git status && exit 1)
	@echo "✅ Git status clean"

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(OUTPUT) coverage.txt