import { useState } from 'react'
import { WebhookList } from '@/components/webhooks/WebhookList'
import { WebhookDetail } from '@/components/webhooks/WebhookDetail'
import { CreateWebhookDialog } from '@/components/webhooks/CreateWebhookDialog'
import type { WebhookEndpoint } from '@/hooks/use-api'

export function WebhooksPage() {
  const [selectedWebhook, setSelectedWebhook] = useState<WebhookEndpoint | null>(null)

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight text-grey-900 dark:text-white">
            Webhooks
          </h1>
          <p className="mt-1 text-sm text-grey-500">
            Manage webhook endpoints and inspect request history
          </p>
        </div>
        <CreateWebhookDialog />
      </div>

      {/* Content */}
      <div className={selectedWebhook ? 'grid gap-6 lg:grid-cols-[1fr_420px]' : ''}>
        <WebhookList
          onSelectWebhook={setSelectedWebhook}
          selectedWebhookId={selectedWebhook?.id}
        />
        {selectedWebhook && (
          <div className="hidden lg:block">
            <div className="sticky top-24 max-h-[calc(100vh-8rem)] overflow-hidden">
              <WebhookDetail
                webhook={selectedWebhook}
                onClose={() => setSelectedWebhook(null)}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
