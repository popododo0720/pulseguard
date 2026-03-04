import { cn } from '@/lib/utils'
import { useRecentActivity, formatRelativeTime } from '@/hooks/use-mock-data'
import type { JobStatus } from '@/lib/mock-data'

function StatusDot({ status }: { status: JobStatus }) {
  return (
    <span
      className={cn(
        'inline-block h-2 w-2 rounded-full',
        status === 'success' && 'bg-green-500',
        status === 'failure' && 'bg-red-500',
        status === 'running' && 'bg-blue-500 animate-pulse-dot',
        status === 'unknown' && 'bg-grey-400',
      )}
    />
  )
}

function StatusLabel({ status }: { status: JobStatus }) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium',
        status === 'success' && 'bg-green-500/10 text-green-500',
        status === 'failure' && 'bg-red-500/10 text-red-500',
        status === 'running' && 'bg-blue-500/10 text-blue-500',
        status === 'unknown' && 'bg-grey-200 text-grey-700 dark:bg-grey-700/30 dark:text-grey-400',
      )}
    >
      <StatusDot status={status} />
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  )
}

export function RecentActivity() {
  const executions = useRecentActivity()

  return (
    <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]" style={{ animationDelay: '250ms' }}>
      <div className="mb-4">
        <h3 className="text-base font-semibold text-grey-900 dark:text-white">
          Recent Activity
        </h3>
        <p className="mt-0.5 text-sm text-grey-500">Latest job executions</p>
      </div>
      <div className="divide-y divide-grey-100 dark:divide-white/5">
        {executions.map((exec) => (
          <div
            key={exec.id}
            className="flex items-center justify-between py-3.5 first:pt-0 last:pb-0"
          >
            <div className="flex items-center gap-3 min-w-0">
              <StatusDot status={exec.status} />
              <div className="min-w-0">
                <p className="truncate text-sm font-medium text-grey-900 dark:text-grey-100">
                  {exec.jobName}
                </p>
                <p className="mt-0.5 text-xs text-grey-500">
                  {exec.agentName} &middot; {exec.duration}
                </p>
              </div>
            </div>
            <div className="flex items-center gap-3 shrink-0 ml-4">
              <StatusLabel status={exec.status} />
              <span className="text-xs text-grey-400 w-16 text-right">
                {formatRelativeTime(exec.startedAt)}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
