import { Clock, CheckCircle2, Server, AlertTriangle } from 'lucide-react'
import { StatCard } from '@/components/dashboard/StatCard'
import { StatusChart } from '@/components/dashboard/StatusChart'
import { RecentActivity } from '@/components/dashboard/RecentActivity'
import { useDashboardStats } from '@/hooks/use-api'

export function DashboardPage() {
  const { data: stats, isLoading } = useDashboardStats()

  return (
    <div className="space-y-8">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-grey-900 dark:text-white">
          Dashboard
        </h1>
        <p className="mt-1 text-sm text-grey-500">
          Monitor your cron jobs and webhooks at a glance
        </p>
      </div>

      {/* Stat cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          label="Total Jobs"
          value={isLoading ? '—' : (stats?.total_jobs ?? 0)}
          icon={Clock}
          delay={0}
        />
        <StatCard
          label="Success Rate"
          value={isLoading ? '—' : (stats?.success_rate ?? 0)}
          suffix="%"
          icon={CheckCircle2}
          delay={50}
        />
        <StatCard
          label="Active Agents"
          value={isLoading ? '—' : (stats?.total_agents ?? 0)}
          icon={Server}
          delay={100}
        />
        <StatCard
          label="Recent Failures"
          value={isLoading ? '—' : (stats?.recent_failures?.length ?? 0)}
          icon={AlertTriangle}
          delay={150}
        />
      </div>

      {/* Chart + Activity */}
      <div className="grid gap-8 lg:grid-cols-[1.2fr_1fr]">
        <StatusChart />
        <RecentActivity />
      </div>
    </div>
  )
}
