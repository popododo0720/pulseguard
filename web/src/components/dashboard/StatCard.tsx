import type { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

interface StatCardProps {
  label: string
  value: string | number
  suffix?: string
  icon: LucideIcon
  trend?: { value: number; isPositive: boolean }
  delay?: number
}

export function StatCard({ label, value, suffix, icon: Icon, trend, delay = 0 }: StatCardProps) {
  return (
    <div
      className={cn(
        'animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]',
      )}
      style={{ animationDelay: `${delay}ms` }}
    >
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-grey-500">{label}</span>
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
          <Icon className="h-[18px] w-[18px] text-grey-400" strokeWidth={1.8} />
        </div>
      </div>
      <div className="mt-3 flex items-baseline gap-1.5">
        <span className="text-3xl font-bold tracking-tight text-grey-900 dark:text-white">
          {value}
        </span>
        {suffix && (
          <span className="text-lg font-medium text-grey-400">{suffix}</span>
        )}
      </div>
      {trend && (
        <div className="mt-2 flex items-center gap-1">
          <span
            className={cn(
              'text-xs font-medium',
              trend.isPositive ? 'text-green-500' : 'text-red-500'
            )}
          >
            {trend.isPositive ? '+' : ''}{trend.value}%
          </span>
          <span className="text-xs text-grey-400">vs last week</span>
        </div>
      )}
    </div>
  )
}
