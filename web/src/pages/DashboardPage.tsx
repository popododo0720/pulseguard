import { Clock, CheckCircle2, Server, Webhook } from 'lucide-react'
import { StatCard } from '@/components/dashboard/StatCard'
import { StatusChart } from '@/components/dashboard/StatusChart'
import { RecentActivity } from '@/components/dashboard/RecentActivity'
import { useDashboardStats } from '@/hooks/use-mock-data'

export function DashboardPage() {
  const stats = useDashboardStats()

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
          value={stats.totalJobs}
          icon={Clock}
          trend={{ value: 12, isPositive: true }}
          delay={0}
        />
        <StatCard
          label="Success Rate"
          value={stats.successRate}
          suffix="%"
          icon={CheckCircle2}
          trend={{ value: 2.1, isPositive: true }}
          delay={50}
        />
        <StatCard
          label="Active Agents"
          value={stats.activeAgents}
          icon={Server}
          delay={100}
        />
        <StatCard
          label="Webhooks"
          value={stats.webhookCount}
          icon={Webhook}
          trend={{ value: 8, isPositive: true }}
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
