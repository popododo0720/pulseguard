.PHONY: all proto build build-server build-agent clean dev dev-web test lint

# Go binary output
BIN_DIR := bin
SERVER_BIN := $(BIN_DIR)/pulseguard-server
AGENT_BIN := $(BIN_DIR)/pulseguard-agent

# Go build flags
GO_FLAGS := -trimpath
LDFLAGS := -s -w

all: proto build

# ============================================================================
# Protobuf
# ============================================================================

proto:
	@echo "==> Generating protobuf code..."
	@PATH="$(PATH):$(shell go env GOPATH)/bin" buf generate

proto-lint:
	@buf lint

# ============================================================================
# Build
# ============================================================================

build: build-server build-agent

build-server:
	@echo "==> Building server..."
	@mkdir -p $(BIN_DIR)
	go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(SERVER_BIN) ./cmd/server

build-agent:
	@echo "==> Building agent..."
	@mkdir -p $(BIN_DIR)
	go build $(GO_FLAGS) -ldflags "$(LDFLAGS)" -o $(AGENT_BIN) ./cmd/agent

# ============================================================================
# Development
# ============================================================================

dev:
	@echo "==> Starting server in dev mode..."
	go run ./cmd/server --dev

dev-web:
	@echo "==> Starting web dev server..."
	cd web && npm run dev

dev-agent:
	@echo "==> Starting agent in dev mode..."
	go run ./cmd/agent --server localhost:9090 --token dev-token

# ============================================================================
# Test & Lint
# ============================================================================

test:
	go test ./... -race -cover

lint:
	golangci-lint run ./...

# ============================================================================
# Web
# ============================================================================

web-install:
	cd web && npm install

web-build:
	cd web && npm run build

# ============================================================================
# Docker
# ============================================================================

docker-build:
	docker build -t pulseguard-server -f Dockerfile .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# ============================================================================
# Clean
# ============================================================================

clean:
	rm -rf $(BIN_DIR)
	rm -rf gen/
	rm -rf web/dist/
