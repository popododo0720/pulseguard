import { cn } from '@/lib/utils'
import { formatRelativeTime } from '@/lib/utils'
import { useAgents } from '@/hooks/use-api'
import { Monitor, Cpu, Clock } from 'lucide-react'

export function AgentList() {
  const { data: agents, isLoading } = useAgents()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-grey-200 border-t-blue-500" />
      </div>
    )
  }

  if (!agents?.length) {
    return (
      <div className="rounded-2xl bg-white p-16 text-center dark:bg-[#161b22]">
        <p className="text-sm text-grey-400">No agents connected yet</p>
      </div>
    )
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {agents.map((agent, i) => (
        <div
          key={agent.id}
          className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]"
          style={{ animationDelay: `${i * 60}ms` }}
        >
          {/* Header */}
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
                <Monitor className="h-5 w-5 text-grey-400" strokeWidth={1.8} />
              </div>
              <div>
                <h3 className="text-sm font-semibold text-grey-900 dark:text-white">
                  {agent.hostname || agent.name}
                </h3>
                <p className="mt-0.5 text-xs text-grey-500">{agent.ip_address}</p>
              </div>
            </div>
            <span
              className={cn(
                'inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium',
                agent.status === 'online'
                  ? 'bg-green-500/10 text-green-500'
                  : 'bg-grey-200 text-grey-500 dark:bg-grey-700/30 dark:text-grey-500',
              )}
            >
              <span
                className={cn(
                  'h-1.5 w-1.5 rounded-full',
                  agent.status === 'online'
                    ? 'bg-green-500 animate-pulse-dot'
                    : 'bg-grey-400',
                )}
              />
              {agent.status === 'online' ? 'Online' : 'Offline'}
            </span>
          </div>

          {/* Info */}
          <div className="mt-5 space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Cpu className="h-3.5 w-3.5 text-grey-400" />
                <span className="text-xs text-grey-500">OS</span>
              </div>
              <span className="text-xs font-medium text-grey-700 dark:text-grey-300">
                {agent.os}{agent.arch ? ` (${agent.arch})` : ''}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Clock className="h-3.5 w-3.5 text-grey-400" />
                <span className="text-xs text-grey-500">Last Heartbeat</span>
              </div>
              <span className="text-xs font-medium text-grey-700 dark:text-grey-300">
                {formatRelativeTime(agent.last_heartbeat_at)}
              </span>
            </div>
          </div>

          {/* Footer */}
          <div className="mt-5 flex items-center justify-end border-t border-grey-100 pt-4 dark:border-white/5">
            <span className="rounded-md bg-grey-100 px-2 py-0.5 text-[10px] text-grey-500 dark:bg-white/5 dark:text-grey-500">
              v{agent.agent_version}
            </span>
          </div>
        </div>
      ))}
    </div>
  )
}
