import { cn } from '@/lib/utils'
import { formatRelativeTime } from '@/lib/utils'
import { useWebhookEndpoints } from '@/hooks/use-api'
import type { WebhookEndpoint } from '@/hooks/use-api'
import { Copy, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface WebhookListProps {
  onSelectWebhook: (webhook: WebhookEndpoint) => void
  selectedWebhookId?: string
}

export function WebhookList({ onSelectWebhook, selectedWebhookId }: WebhookListProps) {
  const { data: webhooks, isLoading } = useWebhookEndpoints()

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-16">
        <div className="h-5 w-5 animate-spin rounded-full border-2 border-grey-200 border-t-blue-500" />
      </div>
    )
  }

  if (!webhooks?.length) {
    return (
      <div className="rounded-2xl bg-white p-16 text-center dark:bg-[#161b22]">
        <p className="text-sm text-grey-400">No webhook endpoints configured yet</p>
      </div>
    )
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2">
      {webhooks.map((wh, i) => (
        <div
          key={wh.id}
          onClick={() => onSelectWebhook(wh)}
          className={cn(
            'animate-fade-in-up cursor-pointer rounded-2xl bg-white p-6 transition-colors duration-150 dark:bg-[#161b22]',
            selectedWebhookId === wh.id
              ? 'ring-2 ring-blue-500/30'
              : 'hover:bg-grey-50/50 dark:hover:bg-white/[0.02]',
          )}
          style={{ animationDelay: `${i * 60}ms` }}
        >
          <div className="flex items-start justify-between">
            <div className="min-w-0 flex-1">
              <div className="flex items-center gap-2">
                <h3 className="truncate text-sm font-semibold text-grey-900 dark:text-white">
                  {wh.name}
                </h3>
                <span
                  className={cn(
                    'inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-medium',
                    wh.enabled
                      ? 'bg-green-500/10 text-green-500'
                      : 'bg-grey-200 text-grey-500 dark:bg-grey-700/30 dark:text-grey-500',
                  )}
                >
                  <span
                    className={cn(
                      'h-1 w-1 rounded-full',
                      wh.enabled ? 'bg-green-500 animate-pulse-dot' : 'bg-grey-400',
                    )}
                  />
                  {wh.enabled ? 'Active' : 'Inactive'}
                </span>
              </div>
              <p className="mt-1 text-xs text-grey-500">{wh.target_url}</p>
            </div>
          </div>

          <div className="mt-4 flex items-center gap-1.5 rounded-xl bg-grey-50 px-3 py-2 dark:bg-black/20">
            <code className="flex-1 truncate text-xs text-grey-600 dark:text-grey-400">
              /wh/{wh.slug}
            </code>
            <Button
              variant="ghost"
              size="icon"
              className="toss-press h-6 w-6 text-grey-400 hover:text-blue-500"
              onClick={(e) => e.stopPropagation()}
            >
              <Copy className="h-3 w-3" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="toss-press h-6 w-6 text-grey-400 hover:text-blue-500"
              onClick={(e) => e.stopPropagation()}
            >
              <ExternalLink className="h-3 w-3" />
            </Button>
          </div>

          <div className="mt-4 flex items-center justify-between">
            <div className="flex items-center gap-4">
              <div>
                <span className="text-lg font-bold text-grey-900 dark:text-white">
                  {(wh.request_count ?? 0).toLocaleString()}
                </span>
                <span className="ml-1 text-xs text-grey-500">requests</span>
              </div>
            </div>
            <span className="text-xs text-grey-400">
              Last: {formatRelativeTime(wh.last_request_at)}
            </span>
          </div>
        </div>
      ))}
    </div>
  )
}
