# ============================================================================
# Stage 1: Build Go binaries
# ============================================================================
FROM golang:1.25-alpine AS go-builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pulseguard-server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/pulseguard-agent ./cmd/agent

# ============================================================================
# Stage 2: Build React frontend
# ============================================================================
FROM node:22-alpine AS web-builder

WORKDIR /build/web

COPY web/package.json web/package-lock.json* ./
RUN npm ci

COPY web/ .
RUN npm run build

# ============================================================================
# Stage 3: Server runtime
# ============================================================================
FROM alpine:3.21 AS server

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=go-builder /out/pulseguard-server /app/pulseguard-server
COPY --from=web-builder /build/web/dist /app/web/dist
COPY migrations/ /app/migrations/

EXPOSE 8080 9090

VOLUME ["/app/data"]

ENV PULSEGUARD_DB_PATH=/app/data/pulseguard.db

ENTRYPOINT ["/app/pulseguard-server"]
CMD ["--port", "8080", "--grpc-port", "9090", "--db", "/app/data/pulseguard.db", "--web-dir", "/app/web/dist"]

# ============================================================================
# Stage 4: Agent runtime
# ============================================================================
FROM alpine:3.21 AS agent

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=go-builder /out/pulseguard-agent /app/pulseguard-agent

ENTRYPOINT ["/app/pulseguard-agent"]
