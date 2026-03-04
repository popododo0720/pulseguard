import { formatDateTime } from '@/lib/utils'
import { useWebhookRequests } from '@/hooks/use-api'
import type { WebhookEndpoint } from '@/hooks/use-api'
import { X } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface WebhookDetailProps {
  webhook: WebhookEndpoint
  onClose: () => void
}

function tryParseJson(str: string): string {
  if (!str) return ''
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

export function WebhookDetail({ webhook, onClose }: WebhookDetailProps) {
  const { data: requests, isLoading } = useWebhookRequests(webhook.id)

  return (
    <div className="animate-fade-in-up flex h-full flex-col rounded-2xl bg-white dark:bg-[#161b22]">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-grey-100 p-6 dark:border-white/5">
        <div>
          <h3 className="text-lg font-semibold text-grey-900 dark:text-white">{webhook.name}</h3>
          <p className="mt-0.5 text-sm text-grey-500">{webhook.target_url}</p>
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

      {/* Request History */}
      <div className="flex-1 overflow-y-auto p-6">
        <h4 className="text-sm font-semibold text-grey-900 dark:text-white">Recent Requests</h4>
        <div className="mt-3 space-y-3">
          {isLoading ? (
            <div className="flex items-center justify-center py-4">
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-grey-200 border-t-blue-500" />
            </div>
          ) : !requests?.length ? (
            <p className="py-4 text-center text-sm text-grey-500">No requests yet</p>
          ) : (
            requests.map((req) => (
              <div
                key={req.id}
                className="rounded-xl border border-grey-100 p-4 dark:border-white/5"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <span className="rounded-md bg-grey-100 px-2 py-0.5 text-[10px] font-semibold uppercase text-grey-600 dark:bg-white/5 dark:text-grey-400">
                      {req.method}
                    </span>
                    {req.source_ip && (
                      <span className="text-xs text-grey-400">{req.source_ip}</span>
                    )}
                  </div>
                  <span className="text-xs text-grey-400">
                    {formatDateTime(req.received_at)}
                  </span>
                </div>

                {req.body && (
                  <div className="mt-3">
                    <span className="text-[10px] font-medium uppercase tracking-wider text-grey-400">
                      Body
                    </span>
                    <pre className="mt-1 overflow-x-auto rounded-lg bg-grey-900 px-3 py-2 text-xs text-grey-300 dark:bg-black/40">
                      {tryParseJson(req.body)}
                    </pre>
                  </div>
                )}

                {req.headers && Object.keys(req.headers).length > 0 && (
                  <div className="mt-2">
                    <span className="text-[10px] font-medium uppercase tracking-wider text-grey-400">
                      Headers
                    </span>
                    <pre className="mt-1 overflow-x-auto rounded-lg bg-grey-50 px-3 py-2 text-xs text-grey-600 dark:bg-white/5 dark:text-grey-400">
                      {JSON.stringify(req.headers, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  )
}
