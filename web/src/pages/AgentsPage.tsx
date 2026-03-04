import { AgentList } from '@/components/agents/AgentList'

export function AgentsPage() {
  return (
    <div className="space-y-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-grey-900 dark:text-white">
          Agents
        </h1>
        <p className="mt-1 text-sm text-grey-500">
          View connected agents and their health status
        </p>
      </div>

      {/* Content */}
      <AgentList />
    </div>
  )
}
