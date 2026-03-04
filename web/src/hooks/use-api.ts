import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

// ─── Types (matching backend snake_case API responses) ────────────────────────

export type JobStatus = 'success' | 'failure' | 'running' | 'unknown'
export type AgentStatus = 'online' | 'offline'
export type WebhookMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'

export interface Agent {
  id: string
  name: string
  hostname: string
  ip_address: string
  os: string
  arch: string
  agent_version: string
  labels: Record<string, string> | null
  status: AgentStatus
  last_heartbeat_at: string
  registered_at: string
}

export interface Job {
  id: string
  agent_id: string
  name: string
  schedule: string
  command: string
  working_dir: string
  timeout_seconds: number
  success_conditions: unknown
  failure_policy: unknown
  env: Record<string, string> | null
  enabled: boolean
  source: string
  last_status: JobStatus
  last_run_at: string
  created_at: string
}

export interface JobExecution {
  id: string
  job_id: string
  agent_id: string
  status: JobStatus
  exit_code: number
  stdout: string
  stderr: string
  error: string
  started_at: string
  finished_at: string
  duration_ms: number
  trigger: string
  retry_count: number
}

export interface WebhookEndpoint {
  id: string
  name: string
  slug: string
  target_url: string
  headers: Record<string, string> | null
  secret: string
  enabled: boolean
  request_count: number
  last_request_at: string
}

export interface WebhookRequest {
  id: string
  endpoint_id: string
  method: WebhookMethod
  headers: Record<string, string> | null
  body: string
  query_params: Record<string, string> | null
  source_ip: string
  received_at: string
}

export interface DashboardStats {
  total_agents: number
  total_jobs: number
  success_rate: number
  recent_failures: number
}

export interface ChartDataPoint {
  date: string
  success: number
  failure: number
}

// ─── Query Hooks ──────────────────────────────────────────────────────────────

export function useDashboardStats() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: () => api.get<DashboardStats>('/dashboard'),
    refetchInterval: 10000,
  })
}

export function useAgents() {
  return useQuery({
    queryKey: ['agents'],
    queryFn: () => api.get<Agent[]>('/agents'),
    refetchInterval: 10000,
  })
}

export function useJobs() {
  return useQuery({
    queryKey: ['jobs'],
    queryFn: () => api.get<Job[]>('/jobs'),
    refetchInterval: 10000,
  })
}

export function useJobExecutions(jobId?: string) {
  return useQuery({
    queryKey: ['job-executions', jobId],
    queryFn: () => api.get<JobExecution[]>(`/jobs/${jobId}/executions`),
    enabled: !!jobId,
    refetchInterval: 10000,
  })
}

export function useWebhookEndpoints() {
  return useQuery({
    queryKey: ['webhook-endpoints'],
    queryFn: () => api.get<WebhookEndpoint[]>('/webhook-endpoints'),
    refetchInterval: 10000,
  })
}

export function useWebhookRequests(webhookId?: string) {
  return useQuery({
    queryKey: ['webhook-requests', webhookId],
    queryFn: () => api.get<WebhookRequest[]>(`/webhook-endpoints/${webhookId}/requests`),
    enabled: !!webhookId,
    refetchInterval: 10000,
  })
}

// ─── Mutation Hooks ───────────────────────────────────────────────────────────

export function useCreateJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<Job>) => api.post<Job>('/jobs', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['jobs'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    },
  })
}

export function useDeleteJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.delete<void>(`/jobs/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['jobs'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
    },
  })
}

export function useRerunJob() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => api.post<void>(`/jobs/${id}/rerun`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['jobs'] })
      queryClient.invalidateQueries({ queryKey: ['job-executions'] })
    },
  })
}

export function useCreateWebhookEndpoint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<WebhookEndpoint>) => api.post<WebhookEndpoint>('/webhook-endpoints', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['webhook-endpoints'] })
    },
  })
}
