import {
  agents,
  jobs,
  jobExecutions,
  webhookEndpoints,
  webhookRequests,
  chartData,
  dashboardStats,
} from '@/lib/mock-data'
import type {
  Agent,
  Job,
  JobExecution,
  WebhookEndpoint,
  WebhookRequest,
  DashboardStats,
  ChartDataPoint,
} from '@/lib/mock-data'

export function useAgents(): Agent[] {
  return agents
}

export function useJobs(): Job[] {
  return jobs
}

export function useJob(id: string): Job | undefined {
  return jobs.find((j) => j.id === id)
}

export function useJobExecutions(jobId?: string): JobExecution[] {
  if (jobId) {
    return jobExecutions.filter((e) => e.jobId === jobId)
  }
  return jobExecutions
}

export function useWebhookEndpoints(): WebhookEndpoint[] {
  return webhookEndpoints
}

export function useWebhookRequests(webhookId?: string): WebhookRequest[] {
  if (webhookId) {
    return webhookRequests.filter((r) => r.webhookId === webhookId)
  }
  return webhookRequests
}

export function useDashboardStats(): DashboardStats {
  return dashboardStats
}

export function useChartData(): ChartDataPoint[] {
  return chartData
}

export function useRecentActivity(): JobExecution[] {
  return jobExecutions.slice(0, 10)
}

export function formatRelativeTime(dateStr: string): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (seconds < 60) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  if (days < 7) return `${days}d ago`
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

export function formatDateTime(dateStr: string): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
