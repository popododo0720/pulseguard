// ─── Types ────────────────────────────────────────────────────────────────────

export type JobStatus = 'success' | 'failure' | 'running' | 'unknown'
export type AgentStatus = 'online' | 'offline'
export type WebhookMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'

export interface Agent {
  id: string
  hostname: string
  ip: string
  os: string
  status: AgentStatus
  lastHeartbeat: string
  jobCount: number
  version: string
}

export interface Job {
  id: string
  name: string
  schedule: string
  agentId: string
  agentName: string
  status: JobStatus
  lastRun: string
  nextRun: string
  duration: string
  command: string
  description: string
  successCount: number
  failureCount: number
}

export interface JobExecution {
  id: string
  jobId: string
  jobName: string
  status: JobStatus
  startedAt: string
  finishedAt: string
  duration: string
  exitCode: number
  output: string
  agentName: string
}

export interface WebhookEndpoint {
  id: string
  name: string
  slug: string
  targetUrl: string
  method: WebhookMethod
  requestCount: number
  lastRequest: string
  isActive: boolean
  description: string
}

export interface WebhookRequest {
  id: string
  webhookId: string
  method: WebhookMethod
  timestamp: string
  statusCode: number
  duration: string
  headers: Record<string, string>
  body: string
  response: string
}

export interface DashboardStats {
  totalJobs: number
  successRate: number
  activeAgents: number
  webhookCount: number
}

export interface ChartDataPoint {
  date: string
  success: number
  failure: number
}

// ─── Mock Data ────────────────────────────────────────────────────────────────

export const agents: Agent[] = [
  {
    id: 'agt-001',
    hostname: 'prod-web-01',
    ip: '10.0.1.12',
    os: 'Ubuntu 22.04 LTS',
    status: 'online',
    lastHeartbeat: '2026-03-04T10:30:00Z',
    jobCount: 4,
    version: '1.2.0',
  },
  {
    id: 'agt-002',
    hostname: 'prod-worker-01',
    ip: '10.0.1.15',
    os: 'Debian 12',
    status: 'online',
    lastHeartbeat: '2026-03-04T10:29:45Z',
    jobCount: 3,
    version: '1.2.0',
  },
  {
    id: 'agt-003',
    hostname: 'staging-01',
    ip: '10.0.2.8',
    os: 'Ubuntu 24.04 LTS',
    status: 'offline',
    lastHeartbeat: '2026-03-03T18:12:00Z',
    jobCount: 1,
    version: '1.1.3',
  },
]

export const jobs: Job[] = [
  {
    id: 'job-001',
    name: 'Database Backup',
    schedule: '0 2 * * *',
    agentId: 'agt-001',
    agentName: 'prod-web-01',
    status: 'success',
    lastRun: '2026-03-04T02:00:00Z',
    nextRun: '2026-03-05T02:00:00Z',
    duration: '4m 32s',
    command: '/opt/scripts/backup-db.sh',
    description: 'Full PostgreSQL database backup to S3',
    successCount: 287,
    failureCount: 3,
  },
  {
    id: 'job-002',
    name: 'Log Rotation',
    schedule: '0 0 * * *',
    agentId: 'agt-001',
    agentName: 'prod-web-01',
    status: 'success',
    lastRun: '2026-03-04T00:00:00Z',
    nextRun: '2026-03-05T00:00:00Z',
    duration: '12s',
    command: 'logrotate /etc/logrotate.d/app',
    description: 'Rotate application logs daily',
    successCount: 364,
    failureCount: 0,
  },
  {
    id: 'job-003',
    name: 'SSL Certificate Check',
    schedule: '0 8 * * 1',
    agentId: 'agt-001',
    agentName: 'prod-web-01',
    status: 'success',
    lastRun: '2026-03-03T08:00:00Z',
    nextRun: '2026-03-10T08:00:00Z',
    duration: '3s',
    command: '/opt/scripts/check-ssl.sh',
    description: 'Verify SSL certificate expiration dates',
    successCount: 52,
    failureCount: 0,
  },
  {
    id: 'job-004',
    name: 'Email Queue Processor',
    schedule: '*/5 * * * *',
    agentId: 'agt-002',
    agentName: 'prod-worker-01',
    status: 'running',
    lastRun: '2026-03-04T10:25:00Z',
    nextRun: '2026-03-04T10:30:00Z',
    duration: '—',
    command: 'python3 /app/process_emails.py',
    description: 'Process queued transactional emails',
    successCount: 8412,
    failureCount: 23,
  },
  {
    id: 'job-005',
    name: 'Cache Warmup',
    schedule: '0 6 * * *',
    agentId: 'agt-002',
    agentName: 'prod-worker-01',
    status: 'failure',
    lastRun: '2026-03-04T06:00:00Z',
    nextRun: '2026-03-05T06:00:00Z',
    duration: '1m 08s',
    command: '/opt/scripts/warm-cache.sh',
    description: 'Pre-warm Redis cache with frequently accessed data',
    successCount: 89,
    failureCount: 12,
  },
  {
    id: 'job-006',
    name: 'Metrics Aggregation',
    schedule: '*/15 * * * *',
    agentId: 'agt-002',
    agentName: 'prod-worker-01',
    status: 'success',
    lastRun: '2026-03-04T10:15:00Z',
    nextRun: '2026-03-04T10:30:00Z',
    duration: '28s',
    command: 'node /app/aggregate-metrics.js',
    description: 'Aggregate and store performance metrics',
    successCount: 2190,
    failureCount: 5,
  },
  {
    id: 'job-007',
    name: 'Stale Session Cleanup',
    schedule: '30 3 * * *',
    agentId: 'agt-001',
    agentName: 'prod-web-01',
    status: 'success',
    lastRun: '2026-03-04T03:30:00Z',
    nextRun: '2026-03-05T03:30:00Z',
    duration: '7s',
    command: '/opt/scripts/cleanup-sessions.sh',
    description: 'Remove expired user sessions from Redis',
    successCount: 365,
    failureCount: 1,
  },
  {
    id: 'job-008',
    name: 'Staging Deploy Check',
    schedule: '0 9 * * 1-5',
    agentId: 'agt-003',
    agentName: 'staging-01',
    status: 'unknown',
    lastRun: '2026-03-03T09:00:00Z',
    nextRun: '2026-03-04T09:00:00Z',
    duration: '45s',
    command: '/opt/scripts/verify-deploy.sh',
    description: 'Verify staging deployment health after rollout',
    successCount: 41,
    failureCount: 8,
  },
]

export const jobExecutions: JobExecution[] = [
  {
    id: 'exec-001',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'running',
    startedAt: '2026-03-04T10:25:00Z',
    finishedAt: '',
    duration: '—',
    exitCode: -1,
    output: 'Processing batch 4821...',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-002',
    jobId: 'job-006',
    jobName: 'Metrics Aggregation',
    status: 'success',
    startedAt: '2026-03-04T10:15:00Z',
    finishedAt: '2026-03-04T10:15:28Z',
    duration: '28s',
    exitCode: 0,
    output: 'Aggregated 12,847 data points across 34 metrics.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-003',
    jobId: 'job-005',
    jobName: 'Cache Warmup',
    status: 'failure',
    startedAt: '2026-03-04T06:00:00Z',
    finishedAt: '2026-03-04T06:01:08Z',
    duration: '1m 08s',
    exitCode: 1,
    output: 'Error: Redis connection refused at 10.0.1.20:6379. Retry limit exceeded.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-004',
    jobId: 'job-007',
    jobName: 'Stale Session Cleanup',
    status: 'success',
    startedAt: '2026-03-04T03:30:00Z',
    finishedAt: '2026-03-04T03:30:07Z',
    duration: '7s',
    exitCode: 0,
    output: 'Cleaned 1,247 expired sessions.',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-005',
    jobId: 'job-001',
    jobName: 'Database Backup',
    status: 'success',
    startedAt: '2026-03-04T02:00:00Z',
    finishedAt: '2026-03-04T02:04:32Z',
    duration: '4m 32s',
    exitCode: 0,
    output: 'Backup complete. Size: 2.4 GB. Uploaded to s3://backups/db/2026-03-04.sql.gz',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-006',
    jobId: 'job-002',
    jobName: 'Log Rotation',
    status: 'success',
    startedAt: '2026-03-04T00:00:00Z',
    finishedAt: '2026-03-04T00:00:12Z',
    duration: '12s',
    exitCode: 0,
    output: 'Rotated 6 log files. Compressed 48 MB → 3.2 MB.',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-007',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'success',
    startedAt: '2026-03-04T10:20:00Z',
    finishedAt: '2026-03-04T10:20:18Z',
    duration: '18s',
    exitCode: 0,
    output: 'Processed 34 emails. 0 failures.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-008',
    jobId: 'job-006',
    jobName: 'Metrics Aggregation',
    status: 'success',
    startedAt: '2026-03-04T10:00:00Z',
    finishedAt: '2026-03-04T10:00:31Z',
    duration: '31s',
    exitCode: 0,
    output: 'Aggregated 11,203 data points across 34 metrics.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-009',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'success',
    startedAt: '2026-03-04T10:15:00Z',
    finishedAt: '2026-03-04T10:15:11Z',
    duration: '11s',
    exitCode: 0,
    output: 'Processed 12 emails. 0 failures.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-010',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'failure',
    startedAt: '2026-03-04T10:10:00Z',
    finishedAt: '2026-03-04T10:10:05Z',
    duration: '5s',
    exitCode: 1,
    output: 'Error: SMTP server timeout. Could not connect to mail.example.com:587.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-011',
    jobId: 'job-008',
    jobName: 'Staging Deploy Check',
    status: 'failure',
    startedAt: '2026-03-03T09:00:00Z',
    finishedAt: '2026-03-03T09:00:45Z',
    duration: '45s',
    exitCode: 2,
    output: 'Health check failed: /api/health returned 503.',
    agentName: 'staging-01',
  },
  {
    id: 'exec-012',
    jobId: 'job-003',
    jobName: 'SSL Certificate Check',
    status: 'success',
    startedAt: '2026-03-03T08:00:00Z',
    finishedAt: '2026-03-03T08:00:03Z',
    duration: '3s',
    exitCode: 0,
    output: 'All 4 certificates valid. Nearest expiry: api.example.com (62 days).',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-013',
    jobId: 'job-001',
    jobName: 'Database Backup',
    status: 'success',
    startedAt: '2026-03-03T02:00:00Z',
    finishedAt: '2026-03-03T02:05:01Z',
    duration: '5m 01s',
    exitCode: 0,
    output: 'Backup complete. Size: 2.3 GB. Uploaded to s3://backups/db/2026-03-03.sql.gz',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-014',
    jobId: 'job-005',
    jobName: 'Cache Warmup',
    status: 'success',
    startedAt: '2026-03-03T06:00:00Z',
    finishedAt: '2026-03-03T06:00:52Z',
    duration: '52s',
    exitCode: 0,
    output: 'Warmed 1,842 cache entries. Hit rate improvement: 34% → 89%.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-015',
    jobId: 'job-002',
    jobName: 'Log Rotation',
    status: 'success',
    startedAt: '2026-03-03T00:00:00Z',
    finishedAt: '2026-03-03T00:00:09Z',
    duration: '9s',
    exitCode: 0,
    output: 'Rotated 6 log files. Compressed 41 MB → 2.8 MB.',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-016',
    jobId: 'job-006',
    jobName: 'Metrics Aggregation',
    status: 'success',
    startedAt: '2026-03-03T23:45:00Z',
    finishedAt: '2026-03-03T23:45:25Z',
    duration: '25s',
    exitCode: 0,
    output: 'Aggregated 10,421 data points across 34 metrics.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-017',
    jobId: 'job-007',
    jobName: 'Stale Session Cleanup',
    status: 'success',
    startedAt: '2026-03-03T03:30:00Z',
    finishedAt: '2026-03-03T03:30:05Z',
    duration: '5s',
    exitCode: 0,
    output: 'Cleaned 982 expired sessions.',
    agentName: 'prod-web-01',
  },
  {
    id: 'exec-018',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'success',
    startedAt: '2026-03-04T10:05:00Z',
    finishedAt: '2026-03-04T10:05:14Z',
    duration: '14s',
    exitCode: 0,
    output: 'Processed 27 emails. 0 failures.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-019',
    jobId: 'job-004',
    jobName: 'Email Queue Processor',
    status: 'success',
    startedAt: '2026-03-04T10:00:00Z',
    finishedAt: '2026-03-04T10:00:09Z',
    duration: '9s',
    exitCode: 0,
    output: 'Processed 8 emails. 0 failures.',
    agentName: 'prod-worker-01',
  },
  {
    id: 'exec-020',
    jobId: 'job-001',
    jobName: 'Database Backup',
    status: 'failure',
    startedAt: '2026-03-02T02:00:00Z',
    finishedAt: '2026-03-02T02:06:12Z',
    duration: '6m 12s',
    exitCode: 1,
    output: 'Error: S3 upload failed. AccessDenied: insufficient permissions on bucket.',
    agentName: 'prod-web-01',
  },
]

export const webhookEndpoints: WebhookEndpoint[] = [
  {
    id: 'wh-001',
    name: 'GitHub Deploy Hook',
    slug: 'github-deploy',
    targetUrl: 'http://localhost:8080/wh/github-deploy',
    method: 'POST',
    requestCount: 142,
    lastRequest: '2026-03-04T09:45:00Z',
    isActive: true,
    description: 'Triggers deployment pipeline on push to main branch',
  },
  {
    id: 'wh-002',
    name: 'Stripe Payment Webhook',
    slug: 'stripe-payments',
    targetUrl: 'http://localhost:8080/wh/stripe-payments',
    method: 'POST',
    requestCount: 3847,
    lastRequest: '2026-03-04T10:28:00Z',
    isActive: true,
    description: 'Receives payment confirmation events from Stripe',
  },
  {
    id: 'wh-003',
    name: 'Health Check Ping',
    slug: 'health-ping',
    targetUrl: 'http://localhost:8080/wh/health-ping',
    method: 'GET',
    requestCount: 8640,
    lastRequest: '2026-03-04T10:30:00Z',
    isActive: true,
    description: 'External uptime monitor ping endpoint',
  },
  {
    id: 'wh-004',
    name: 'Slack Alert Receiver',
    slug: 'slack-alerts',
    targetUrl: 'http://localhost:8080/wh/slack-alerts',
    method: 'POST',
    requestCount: 56,
    lastRequest: '2026-03-03T15:20:00Z',
    isActive: false,
    description: 'Receives alert notifications from Slack integration',
  },
]

export const webhookRequests: WebhookRequest[] = [
  {
    id: 'req-001',
    webhookId: 'wh-002',
    method: 'POST',
    timestamp: '2026-03-04T10:28:00Z',
    statusCode: 200,
    duration: '45ms',
    headers: { 'Content-Type': 'application/json', 'Stripe-Signature': 'whsec_...' },
    body: '{"type":"payment_intent.succeeded","data":{"amount":4999}}',
    response: '{"received":true}',
  },
  {
    id: 'req-002',
    webhookId: 'wh-003',
    method: 'GET',
    timestamp: '2026-03-04T10:30:00Z',
    statusCode: 200,
    duration: '2ms',
    headers: { 'User-Agent': 'UptimeRobot/2.0' },
    body: '',
    response: '{"status":"ok","uptime":"99.97%"}',
  },
  {
    id: 'req-003',
    webhookId: 'wh-001',
    method: 'POST',
    timestamp: '2026-03-04T09:45:00Z',
    statusCode: 200,
    duration: '120ms',
    headers: { 'Content-Type': 'application/json', 'X-GitHub-Event': 'push' },
    body: '{"ref":"refs/heads/main","commits":[{"message":"fix: resolve auth bug"}]}',
    response: '{"deployment_id":"dep-4821","status":"triggered"}',
  },
  {
    id: 'req-004',
    webhookId: 'wh-002',
    method: 'POST',
    timestamp: '2026-03-04T10:15:00Z',
    statusCode: 200,
    duration: '38ms',
    headers: { 'Content-Type': 'application/json', 'Stripe-Signature': 'whsec_...' },
    body: '{"type":"charge.refunded","data":{"amount":2500}}',
    response: '{"received":true}',
  },
  {
    id: 'req-005',
    webhookId: 'wh-001',
    method: 'POST',
    timestamp: '2026-03-03T14:22:00Z',
    statusCode: 500,
    duration: '3012ms',
    headers: { 'Content-Type': 'application/json', 'X-GitHub-Event': 'push' },
    body: '{"ref":"refs/heads/main","commits":[{"message":"feat: add dashboard"}]}',
    response: '{"error":"deployment pipeline timeout"}',
  },
  {
    id: 'req-006',
    webhookId: 'wh-004',
    method: 'POST',
    timestamp: '2026-03-03T15:20:00Z',
    statusCode: 200,
    duration: '28ms',
    headers: { 'Content-Type': 'application/json' },
    body: '{"text":"Alert: CPU usage above 90% on prod-web-01"}',
    response: '{"ok":true}',
  },
]

export const chartData: ChartDataPoint[] = [
  { date: 'Feb 26', success: 42, failure: 3 },
  { date: 'Feb 27', success: 45, failure: 1 },
  { date: 'Feb 28', success: 48, failure: 2 },
  { date: 'Mar 01', success: 44, failure: 4 },
  { date: 'Mar 02', success: 46, failure: 2 },
  { date: 'Mar 03', success: 50, failure: 3 },
  { date: 'Mar 04', success: 38, failure: 2 },
]

export const dashboardStats: DashboardStats = {
  totalJobs: 8,
  successRate: 96.4,
  activeAgents: 2,
  webhookCount: 4,
}
