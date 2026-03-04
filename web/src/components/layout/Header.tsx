import { Moon, Sun, Bell } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
import { useDashboardStats } from '@/hooks/use-api'

export function Header() {
  const [isDark, setIsDark] = useState(false)
  const [showNotifications, setShowNotifications] = useState(false)
  const notifRef = useRef<HTMLDivElement>(null)
  const { data: dashboard } = useDashboardStats()

  useEffect(() => {
    const root = document.documentElement
    if (isDark) {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }, [isDark])

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (notifRef.current && !notifRef.current.contains(e.target as Node)) {
        setShowNotifications(false)
      }
    }
    if (showNotifications) {
      document.addEventListener('mousedown', handleClickOutside)
    }
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [showNotifications])

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-grey-200 bg-grey-50/80 px-8 backdrop-blur-sm dark:border-white/[0.06] dark:bg-[#0d1117]/80">
      {/* Brand */}
      <div className="flex items-center gap-2">
        <span className="text-base font-semibold text-grey-900 dark:text-white tracking-tight">PulseGuard</span>
      </div>

      {/* Actions */}
      <div className="flex items-center gap-1">
        <div className="relative" ref={notifRef}>
          <Button
            variant="ghost"
            size="icon"
            className="toss-press h-9 w-9 text-grey-500 hover:bg-grey-100 hover:text-grey-700 dark:hover:bg-white/5"
            onClick={() => setShowNotifications(!showNotifications)}
          >
            <Bell className="h-[18px] w-[18px]" strokeWidth={1.8} />
            {dashboard?.recent_failures && dashboard.recent_failures.length > 0 && (
              <span className="absolute -top-0.5 -right-0.5 h-2 w-2 rounded-full bg-red-500" />
            )}
          </Button>
          {showNotifications && (
            <div className="absolute right-0 top-full mt-2 w-80 rounded-xl border border-grey-200 bg-white p-2 shadow-lg dark:border-white/10 dark:bg-[#161b22]">
              <div className="px-3 py-2 text-xs font-medium text-grey-500 uppercase tracking-wider">Recent Failures</div>
              {(!dashboard?.recent_failures || dashboard.recent_failures.length === 0) ? (
                <div className="px-3 py-4 text-center text-sm text-grey-400">No recent failures</div>
              ) : (
                <div className="max-h-64 overflow-y-auto space-y-0.5">
                  {dashboard.recent_failures.slice(0, 8).map((f: any, i: number) => (
                    <div key={i} className="flex items-center gap-3 rounded-lg px-3 py-2 hover:bg-grey-50 dark:hover:bg-white/5 cursor-default">
                      <div className="h-2 w-2 shrink-0 rounded-full bg-red-500" />
                      <div className="min-w-0 flex-1">
                        <div className="truncate text-sm font-medium text-grey-700 dark:text-grey-300">
                          {f.job_name || f.name || 'Unknown job'}
                        </div>
                        <div className="text-xs text-grey-400">
                          {f.started_at ? new Date(f.started_at).toLocaleString() : 'recently'}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
        <Button
          variant="ghost"
          size="icon"
          className="toss-press h-9 w-9 text-grey-500 hover:bg-grey-100 hover:text-grey-700 dark:hover:bg-white/5"
          onClick={() => setIsDark(!isDark)}
        >
          {isDark ? (
            <Sun className="h-[18px] w-[18px]" strokeWidth={1.8} />
          ) : (
            <Moon className="h-[18px] w-[18px]" strokeWidth={1.8} />
          )}
        </Button>
      </div>
    </header>
  )
}
