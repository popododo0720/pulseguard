import { Bell, Key, Globe, Shield } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'

export function SettingsPage() {
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
              defaultValue="http://localhost:8080"
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
              defaultValue="30"
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
            <Switch />
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
              pg_tok_••••••••••••••••••••a4b7
            </code>
            <Button
              variant="ghost"
              className="toss-press h-8 px-3 text-xs text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-500/10"
            >
              Copy
            </Button>
          </div>
          <Button
            variant="ghost"
            className="toss-press h-9 rounded-xl border border-grey-200 px-4 text-sm text-grey-700 hover:bg-grey-50 dark:border-white/10 dark:text-grey-300 dark:hover:bg-white/5"
          >
            Regenerate Token
          </Button>
        </div>
      </div>

      {/* Security */}
      <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]" style={{ animationDelay: '180ms' }}>
        <div className="flex items-center gap-3 mb-6">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-grey-50 dark:bg-white/5">
            <Shield className="h-[18px] w-[18px] text-grey-400" strokeWidth={1.8} />
          </div>
          <div>
            <h2 className="text-base font-semibold text-grey-900 dark:text-white">Security</h2>
            <p className="text-sm text-grey-500">Webhook signature verification</p>
          </div>
        </div>

        <div className="space-y-5">
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Require HMAC signature
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Validate webhook request signatures</p>
            </div>
            <Switch defaultChecked />
          </div>
          <Separator className="bg-grey-100 dark:bg-white/5" />
          <div className="flex items-center justify-between">
            <div>
              <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
                Rate limiting
              </Label>
              <p className="text-xs text-grey-500 mt-0.5">Limit webhook requests per minute</p>
            </div>
            <Input
              defaultValue="100"
              className="h-9 w-24 rounded-xl border-grey-200 bg-white text-sm text-right dark:border-white/10 dark:bg-white/5"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
