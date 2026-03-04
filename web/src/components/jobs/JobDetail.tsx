import { cn } from '@/lib/utils'
import { useJobExecutions, formatDateTime } from '@/hooks/use-mock-data'
import { X, Play, Terminal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import type { Job, JobStatus } from '@/lib/mock-data'

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

interface JobDetailProps {
  job: Job
  onClose: () => void
}

export function JobDetail({ job, onClose }: JobDetailProps) {
  const executions = useJobExecutions(job.id)

  return (
    <div className="animate-fade-in-up flex h-full flex-col rounded-2xl bg-white dark:bg-[#161b22]">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-grey-100 p-6 dark:border-white/5">
        <div>
          <h3 className="text-lg font-semibold text-grey-900 dark:text-white">{job.name}</h3>
          <p className="mt-0.5 text-sm text-grey-500">{job.description}</p>
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="toss-press h-8 w-8 text-grey-400"
          onClick={onClose}
        >
          <X className="h-4 w-4" />
        </Button>
      </div>

      {/* Info */}
      <div className="space-y-4 p-6">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <span className="text-xs text-grey-500">Status</span>
            <div className="mt-1">
              <StatusBadge status={job.status} />
            </div>
          </div>
          <div>
            <span className="text-xs text-grey-500">Agent</span>
            <p className="mt-1 text-sm font-medium text-grey-900 dark:text-grey-100">{job.agentName}</p>
          </div>
          <div>
            <span className="text-xs text-grey-500">Schedule</span>
            <code className="mt-1 block rounded-md bg-grey-100 px-2 py-0.5 text-xs text-grey-600 dark:bg-white/5 dark:text-grey-400 w-fit">
              {job.schedule}
            </code>
          </div>
          <div>
            <span className="text-xs text-grey-500">Duration</span>
            <p className="mt-1 text-sm font-medium text-grey-900 dark:text-grey-100">{job.duration}</p>
          </div>
        </div>

        <div>
          <span className="text-xs text-grey-500">Command</span>
          <div className="mt-1 flex items-center gap-2 rounded-xl bg-grey-900 px-3 py-2.5 dark:bg-black/40">
            <Terminal className="h-3.5 w-3.5 text-grey-500" />
            <code className="text-xs text-grey-300">{job.command}</code>
          </div>
        </div>

        <div className="flex gap-3 pt-1">
          <Button
            className="toss-press h-9 rounded-xl bg-blue-500 px-4 text-sm font-medium text-white hover:bg-blue-600"
          >
            <Play className="mr-1.5 h-3.5 w-3.5" />
            Run Now
          </Button>
        </div>
      </div>

      <Separator className="bg-grey-100 dark:bg-white/5" />

      {/* Execution History */}
      <div className="flex-1 overflow-y-auto p-6">
        <h4 className="text-sm font-semibold text-grey-900 dark:text-white">Execution History</h4>
        <div className="mt-3 space-y-2">
          {executions.length === 0 ? (
            <p className="py-4 text-center text-sm text-grey-500">No executions yet</p>
          ) : (
            executions.map((exec) => (
              <div
                key={exec.id}
                className="rounded-xl border border-grey-100 p-3 dark:border-white/5"
              >
                <div className="flex items-center justify-between">
                  <StatusBadge status={exec.status} />
                  <span className="text-xs text-grey-400">{formatDateTime(exec.startedAt)}</span>
                </div>
                <p className="mt-2 truncate text-xs text-grey-600 dark:text-grey-400">
                  {exec.output}
                </p>
                <div className="mt-1 flex items-center gap-3 text-xs text-grey-400">
                  <span>Duration: {exec.duration}</span>
                  {exec.exitCode >= 0 && <span>Exit: {exec.exitCode}</span>}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
