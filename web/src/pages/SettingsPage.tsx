import { useState, useEffect } from 'react'
import { Bell, Key, Globe } from 'lucide-react'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'

export function SettingsPage() {
  const [refreshInterval, setRefreshInterval] = useState(() =>
    localStorage.getItem('pulseguard_refresh_interval') || '30'
  )
  const [darkMode, setDarkMode] = useState(() =>
    localStorage.getItem('pulseguard_dark_mode') === 'true'
  )

  const { data: serverSettings } = useQuery({
    queryKey: ['settings'],
    queryFn: () => api.get<{ server_version: string; token_masked: string }>('/settings'),
  })

  useEffect(() => {
    localStorage.setItem('pulseguard_dark_mode', String(darkMode))
    if (darkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }, [darkMode])

  useEffect(() => {
    localStorage.setItem('pulseguard_refresh_interval', refreshInterval)
  }, [refreshInterval])

  return (
    <div className="space-y-8">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-semibold tracking-tight text-grey-900 dark:text-white">
          Settings
        </h1>
        <p className="mt-1 text-sm text-grey-500">
          Configure notifications, API access, and general preferences
        </p>
      </div>

      {/* General Settings */}
      <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]">
        <div className="flex items-center gap-3 mb-6">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
            <Globe className="h-[18px] w-[18px] text-grey-400" strokeWidth={1.8} />
          </div>
          <div>
            <h2 className="text-base font-semibold text-grey-900 dark:text-white">General</h2>
            <p className="text-sm text-grey-500">Basic application settings</p>
          </div>
        </div>

        <div className="space-y-5">
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Server URL
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">The base URL of your PulseGuard server</p>
            </div>
            <Input
              value={window.location.origin}
              readOnly
              className="h-9 w-64 rounded-xl border-grey-200 bg-white text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <Separator className="bg-grey-100 dark:bg-white/5" />
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Auto-refresh interval
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">How often to poll for updates</p>
            </div>
            <Input
              value={refreshInterval}
              onChange={(e) => setRefreshInterval(e.target.value)}
              className="h-9 w-24 rounded-xl border-grey-200 bg-white text-sm text-right dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <Separator className="bg-grey-100 dark:bg-white/5" />
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Dark mode
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Toggle dark theme</p>
            </div>
            <Switch checked={darkMode} onCheckedChange={setDarkMode} />
          </div>
        </div>
      </div>

      {/* Notification Settings */}
      <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]" style={{ animationDelay: '60ms' }}>
        <div className="flex items-center gap-3 mb-6">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
            <Bell className="h-[18px] w-[18px] text-grey-400" strokeWidth={1.8} />
          </div>
          <div>
            <h2 className="text-base font-semibold text-grey-900 dark:text-white">Notifications</h2>
            <p className="text-sm text-grey-500">Configure alert channels</p>
          </div>
        </div>

        <div className="space-y-5">
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Email notifications
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Receive alerts via email</p>
            </div>
            <Switch defaultChecked />
          </div>
          <Separator className="bg-grey-100 dark:bg-white/5" />
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Slack integration
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Send alerts to Slack channels</p>
            </div>
            <Switch />
          </div>
          <Separator className="bg-grey-100 dark:bg-white/5" />
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Notify on failure only
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Only send alerts when jobs fail</p>
            </div>
            <Switch defaultChecked />
          </div>
        </div>
      </div>

      {/* API Token */}
      <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]" style={{ animationDelay: '120ms' }}>
        <div className="flex items-center gap-3 mb-6">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
            <Key className="h-[18px] w-[18px] text-grey-400" strokeWidth={1.8} />
          </div>
          <div>
            <h2 className="text-base font-semibold text-grey-900 dark:text-white">API Token</h2>
            <p className="text-sm text-grey-500">Access token for the PulseGuard API</p>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center gap-2 rounded-xl bg-grey-50 px-4 py-3 dark:bg-black/20">
            <code className="flex-1 text-sm text-grey-600 dark:text-grey-400">
              {serverSettings?.token_masked || '\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022'}
            </code>
            <Button
              variant="ghost"
              className="toss-press h-8 px-3 text-xs text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-500/10"
              onClick={() => {
                if (serverSettings?.token_masked) {
                  navigator.clipboard.writeText(serverSettings.token_masked)
                }
              }}
            >
              Copy
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
