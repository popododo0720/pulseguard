-- PulseGuard Initial Schema
-- SQLite with WAL mode for concurrent read/write

PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

-- ============================================================================
-- Agents
-- ============================================================================

CREATE TABLE IF NOT EXISTS agents (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',
    hostname TEXT NOT NULL DEFAULT '',
    ip_address TEXT NOT NULL DEFAULT '',
    os TEXT NOT NULL DEFAULT '',
    arch TEXT NOT NULL DEFAULT '',
    agent_version TEXT NOT NULL DEFAULT '',
    labels TEXT NOT NULL DEFAULT '{}',
    status TEXT NOT NULL DEFAULT 'offline',
    last_heartbeat_at DATETIME,
    registered_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================================
-- Jobs
-- ============================================================================

CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    schedule TEXT NOT NULL,
    command TEXT NOT NULL,
    working_dir TEXT NOT NULL DEFAULT '',
    timeout_seconds INTEGER NOT NULL DEFAULT 3600,
    success_conditions TEXT NOT NULL DEFAULT '{"expected_exit_code": 0}',
    failure_policy TEXT NOT NULL DEFAULT '{"max_retries": 0, "retry_delay_seconds": 60}',
    env TEXT NOT NULL DEFAULT '{}',
    enabled INTEGER NOT NULL DEFAULT 1,
    source TEXT NOT NULL DEFAULT 'manual',
    last_status TEXT NOT NULL DEFAULT 'unknown',
    last_run_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_jobs_agent_id ON jobs(agent_id);
CREATE INDEX IF NOT EXISTS idx_jobs_enabled ON jobs(enabled);
CREATE INDEX IF NOT EXISTS idx_jobs_last_status ON jobs(last_status);

-- ============================================================================
-- Job Executions
-- ============================================================================

CREATE TABLE IF NOT EXISTS job_executions (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'running',
    exit_code INTEGER,
    stdout TEXT NOT NULL DEFAULT '',
    stderr TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    started_at DATETIME NOT NULL,
    finished_at DATETIME,
    duration_ms INTEGER,
    trigger TEXT NOT NULL DEFAULT 'scheduled',
    retry_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_job_executions_job_id ON job_executions(job_id);
CREATE INDEX IF NOT EXISTS idx_job_executions_agent_id ON job_executions(agent_id);
CREATE INDEX IF NOT EXISTS idx_job_executions_status ON job_executions(status);
CREATE INDEX IF NOT EXISTS idx_job_executions_started_at ON job_executions(started_at DESC);

-- ============================================================================
-- Webhook Endpoints
-- ============================================================================

CREATE TABLE IF NOT EXISTS webhook_endpoints (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    target_url TEXT NOT NULL,
    headers TEXT NOT NULL DEFAULT '{}',
    secret TEXT NOT NULL DEFAULT '',
    enabled INTEGER NOT NULL DEFAULT 1,
    request_count INTEGER NOT NULL DEFAULT 0,
    last_request_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================================
-- Webhook Requests (received)
-- ============================================================================

CREATE TABLE IF NOT EXISTS webhook_requests (
    id TEXT PRIMARY KEY,
    endpoint_id TEXT NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    method TEXT NOT NULL,
    headers TEXT NOT NULL DEFAULT '{}',
    body TEXT NOT NULL DEFAULT '',
    query_params TEXT NOT NULL DEFAULT '',
    source_ip TEXT NOT NULL DEFAULT '',
    received_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_requests_endpoint_id ON webhook_requests(endpoint_id);
CREATE INDEX IF NOT EXISTS idx_webhook_requests_received_at ON webhook_requests(received_at DESC);

-- ============================================================================
-- Webhook Deliveries (forwarded to target)
-- ============================================================================

CREATE TABLE IF NOT EXISTS webhook_deliveries (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL REFERENCES webhook_requests(id) ON DELETE CASCADE,
    endpoint_id TEXT NOT NULL REFERENCES webhook_endpoints(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    response_status INTEGER,
    response_body TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    attempt_number INTEGER NOT NULL DEFAULT 1,
    delivered_at DATETIME,
    next_retry_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_request_id ON webhook_deliveries(request_id);
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_status ON webhook_deliveries(status);
CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_next_retry ON webhook_deliveries(next_retry_at)
    WHERE status = 'pending' AND next_retry_at IS NOT NULL;

-- ============================================================================
-- Notification Channels
-- ============================================================================

CREATE TABLE IF NOT EXISTS notification_channels (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    config TEXT NOT NULL DEFAULT '{}',
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);
