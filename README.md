# PulseGuard

Open-source, self-hosted cron job monitoring + webhook monitoring/replay tool.

Monitor cron jobs across multiple servers, receive and replay webhooks, get notified on failures — all from a single dashboard.

## Features

- **Cron Job Monitoring** — Track job executions, success/failure status, stdout/stderr capture
- **One-Click Rerun** — Re-execute failed jobs directly from the web UI
- **Webhook Proxy** — Receive webhooks, store payloads, forward to targets, auto-retry on failure
- **Webhook Replay** — Replay any received webhook with one click
- **Multi-Server** — Deploy agents across multiple servers, monitor everything centrally
- **Crontab Discovery** — Agent auto-discovers existing crontab entries
- **Custom Success Conditions** — Define success by exit code, stdout content, timeout, and more
- **Real-time Updates** — gRPC streaming for instant command delivery to agents

## Architecture

```
┌─────────────┐     gRPC      ┌──────────────────┐     REST      ┌─────────┐
│   Agent #1  │◄─────────────►│                  │◄─────────────►│         │
│  (server-1) │               │   PulseGuard     │               │  Web UI │
├─────────────┤               │     Server       │               │  (SPA)  │
│   Agent #2  │◄─────────────►│                  │               │         │
│  (server-2) │               │  REST + gRPC     │               └─────────┘
├─────────────┤               │  SQLite          │
│   Agent #N  │◄─────────────►│  Webhook Proxy   │
│  (server-N) │               └──────────────────┘
└─────────────┘
```

## Quick Start

### Docker Compose (Recommended)

```bash
docker compose up -d
```

Access the dashboard at `http://localhost:8080`.

### From Source

```bash
# Build
make all

# Run server
./bin/pulseguard-server --port 8080 --grpc-port 9090

# Run agent (on each monitored server)
./bin/pulseguard-agent --server your-server:9090 --token your-token
```

### Development

```bash
# Start server in dev mode
make dev

# Start frontend dev server (separate terminal)
make dev-web

# Start agent pointing to local server
make dev-agent
```

## Tech Stack

| Component | Technology |
|-----------|------------|
| Server | Go, gRPC, Chi (REST), SQLite |
| Agent | Go, gRPC client, cron scheduler |
| Frontend | React, Vite, Tailwind CSS v4, shadcn/ui |
| Protocol | Protocol Buffers (gRPC + server streaming) |
| Database | SQLite (WAL mode) |
| Deployment | Docker, multi-stage build |

## Project Structure

```
pulseguard/
├── cmd/
│   ├── server/          # Server entrypoint
│   └── agent/           # Agent entrypoint
├── internal/
│   ├── server/
│   │   ├── api/         # REST API handlers
│   │   ├── grpc/        # gRPC service implementation
│   │   ├── store/       # SQLite data access
│   │   └── webhook/     # Webhook proxy engine
│   ├── agent/
│   │   ├── client/      # gRPC client
│   │   ├── executor/    # Job runner
│   │   ├── discovery/   # Crontab auto-discovery
│   │   └── scheduler/   # Cron scheduler
│   └── models/          # Shared domain models
├── proto/               # Protobuf definitions
├── gen/                 # Generated gRPC code
├── web/                 # React SPA
├── migrations/          # SQL migrations
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Configuration

### Server

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | REST API + Web UI port |
| `--grpc-port` | `9090` | gRPC port for agents |
| `--db` | `./pulseguard.db` | SQLite database path |
| `--dev` | `false` | Enable dev mode (permissive CORS) |

### Agent

| Flag | Default | Description |
|------|---------|-------------|
| `--server` | (required) | Server gRPC address (host:port) |
| `--token` | (required) | Authentication token |

## API Endpoints

### REST API (Web UI)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/dashboard` | Dashboard summary stats |
| GET | `/api/agents` | List agents |
| GET | `/api/jobs` | List jobs |
| POST | `/api/jobs` | Create job |
| POST | `/api/jobs/:id/rerun` | Rerun a job |
| GET | `/api/jobs/:id/executions` | Job execution history |
| GET | `/api/webhook-endpoints` | List webhook endpoints |
| POST | `/api/webhook-endpoints` | Create webhook endpoint |
| POST | `/wh/:slug` | Receive incoming webhook |

### gRPC (Agent ↔ Server)

| RPC | Description |
|-----|-------------|
| `Register` | Agent registration |
| `Heartbeat` | Periodic health check |
| `ReportJobResult` | Report execution result |
| `CommandStream` | Server-streaming commands |
| `ReportDiscoveredJobs` | Report crontab discoveries |

## License

MIT
