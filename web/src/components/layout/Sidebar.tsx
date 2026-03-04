import { NavLink } from 'react-router-dom'
import {
  LayoutDashboard,
  Clock,
  Webhook,
  Server,
  Settings,
  Activity,
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/jobs', icon: Clock, label: 'Jobs' },
  { to: '/webhooks', icon: Webhook, label: 'Webhooks' },
  { to: '/agents', icon: Server, label: 'Agents' },
  { to: '/settings', icon: Settings, label: 'Settings' },
]

export function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 z-40 flex h-screen w-[240px] flex-col bg-grey-900">
      {/* Logo */}
      <div className="flex h-16 items-center gap-2.5 px-6">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-500">
          <Activity className="h-4.5 w-4.5 text-white" strokeWidth={2.5} />
        </div>
        <span className="text-lg font-semibold tracking-tight text-white">
          PulseGuard
        </span>
      </div>

      {/* Navigation */}
      <nav className="mt-2 flex flex-1 flex-col gap-1 px-3">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.to === '/'}
            className={({ isActive }) =>
              cn(
                'toss-press flex items-center gap-3 rounded-xl px-3 py-2.5 text-[15px] font-medium transition-colors duration-150',
                isActive
                  ? 'bg-white/10 text-white'
                  : 'text-grey-400 hover:bg-white/[0.06] hover:text-grey-200'
              )
            }
          >
            <item.icon className="h-[18px] w-[18px]" strokeWidth={1.8} />
            {item.label}
          </NavLink>
        ))}
      </nav>

      {/* Footer */}
      <div className="border-t border-white/[0.06] px-6 py-4">
        <p className="text-xs text-grey-500">PulseGuard v1.0.0</p>
        <p className="mt-0.5 text-xs text-grey-600">Connected to localhost</p>
      </div>
    </aside>
  )
}
