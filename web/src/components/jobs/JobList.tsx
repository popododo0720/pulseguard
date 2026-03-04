import { cn } from '@/lib/utils'
import { formatRelativeTime } from '@/lib/utils'
import { useJobs, useAgents, useRerunJob, useDeleteJob } from '@/hooks/use-api'
import type { Job, JobStatus } from '@/hooks/use-api'
import { useMemo } from 'react'
import { Play, MoreHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

function StatusBadge({ status }: { status: JobStatus }) {
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
      <span
        className={cn(
          'h-1.5 w-1.5 rounded-full',
          status === 'success' && 'bg-green-500',
          status === 'failure' && 'bg-red-500',
          status === 'running' && 'bg-blue-500 animate-pulse-dot',
          status === 'unknown' && 'bg-grey-400',
        )}
      />
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  )
}

interface JobListProps {
  onSelectJob: (job: Job) => void
  selectedJobId?: string
}

export function JobList({ onSelectJob, selectedJobId }: JobListProps) {
  const { data: jobs, isLoading } = useJobs()
  const { data: agents } = useAgents()
  const rerunJob = useRerunJob()
  const deleteJob = useDeleteJob()

  const agentMap = useMemo(() => {
    const map: Record<string, string> = {}
    agents?.forEach((a) => { map[a.id] = a.hostname || a.name })
    return map
  }, [agents])

  if (isLoading) {
    return (
      <div className="animate-fade-in-up rounded-2xl bg-white dark:bg-[#161b22]">
        <div className="flex items-center justify-center py-16">
          <div className="h-5 w-5 animate-spin rounded-full border-2 border-grey-200 border-t-blue-500" />
        </div>
      </div>
    )
  }

  if (!jobs?.length) {
    return (
      <div className="animate-fade-in-up rounded-2xl bg-white dark:bg-[#161b22]">
        <p className="py-16 text-center text-sm text-grey-400">No jobs configured yet</p>
      </div>
    )
  }

  return (
    <div className="animate-fade-in-up rounded-2xl bg-white dark:bg-[#161b22]">
      {/* Table header */}
      <div className="grid grid-cols-[2fr_1fr_1fr_1fr_1fr_auto] gap-4 border-b border-grey-100 px-6 py-3 dark:border-white/5">
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500">Name</span>
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500">Agent</span>
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500">Schedule</span>
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500">Status</span>
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500">Last Run</span>
        <span className="text-xs font-medium uppercase tracking-wider text-grey-500 w-20">Actions</span>
      </div>

      {/* Table rows */}
      {jobs.map((job, i) => (
        <div
          key={job.id}
          onClick={() => onSelectJob(job)}
          className={cn(
            'grid cursor-pointer grid-cols-[2fr_1fr_1fr_1fr_1fr_auto] gap-4 border-b border-grey-50 px-6 py-4 transition-colors duration-150 last:border-b-0 dark:border-white/[0.03]',
            selectedJobId === job.id
              ? 'bg-blue-50/50 dark:bg-blue-500/5'
              : 'hover:bg-grey-50 dark:hover:bg-white/[0.02]',
            'animate-fade-in-up',
          )}
          style={{ animationDelay: `${i * 30}ms` }}
        >
          <div className="min-w-0">
            <p className="truncate text-sm font-medium text-grey-900 dark:text-grey-100">
              {job.name}
            </p>
            <p className="mt-0.5 truncate text-xs text-grey-500">{job.command}</p>
          </div>
          <div className="flex items-center">
            <span className="text-sm text-grey-700 dark:text-grey-300">
              {agentMap[job.agent_id] || job.agent_id || '—'}
            </span>
          </div>
          <div className="flex items-center">
            <code className="rounded-md bg-grey-100 px-2 py-0.5 text-xs text-grey-600 dark:bg-white/5 dark:text-grey-400">
              {job.schedule}
            </code>
          </div>
          <div className="flex items-center">
            <StatusBadge status={job.last_status} />
          </div>
          <div className="flex items-center">
            <span className="text-sm text-grey-500">{formatRelativeTime(job.last_run_at)}</span>
          </div>
          <div className="flex items-center gap-1 w-20 justify-end">
            <Button
              variant="ghost"
              size="icon"
              className="toss-press h-8 w-8 text-grey-400 hover:bg-grey-100 hover:text-blue-500 dark:hover:bg-white/5"
              onClick={(e) => {
                e.stopPropagation()
                rerunJob.mutate(job.id)
              }}
            >
              <Play className="h-3.5 w-3.5" />
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="toss-press h-8 w-8 text-grey-400 hover:bg-grey-100 dark:hover:bg-white/5"
                  onClick={(e) => e.stopPropagation()}
                >
                  <MoreHorizontal className="h-3.5 w-3.5" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-40">
                <DropdownMenuItem>Edit</DropdownMenuItem>
                <DropdownMenuItem>View History</DropdownMenuItem>
                <DropdownMenuItem>Pause</DropdownMenuItem>
                <DropdownMenuItem
                  className="text-red-500"
                  onClick={(e) => {
                    e.stopPropagation()
                    deleteJob.mutate(job.id)
                  }}
                >
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      ))}
    </div>
  )
}
