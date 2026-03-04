import { useState } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Plus } from 'lucide-react'
import { useCreateWebhookEndpoint } from '@/hooks/use-api'

export function CreateWebhookDialog() {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')
  const [slug, setSlug] = useState('')
  const [targetUrl, setTargetUrl] = useState('')
  const createWebhook = useCreateWebhookEndpoint()

  const handleSubmit = () => {
    if (!name || !slug) return
    createWebhook.mutate(
      { name, slug, target_url: targetUrl || undefined } as never,
      {
        onSuccess: () => {
          setOpen(false)
          setName('')
          setSlug('')
          setTargetUrl('')
        },
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="toss-press h-9 rounded-xl bg-blue-500 px-4 text-sm font-medium text-white hover:bg-blue-600">
          <Plus className="mr-1.5 h-4 w-4" />
          New Webhook
        </Button>
      </DialogTrigger>
      <DialogContent className="rounded-2xl border-grey-200 bg-white p-0 sm:max-w-lg dark:border-white/10 dark:bg-[#161b22]">
        <DialogHeader className="border-b border-grey-100 p-6 dark:border-white/5">
          <DialogTitle className="text-lg font-semibold text-grey-900 dark:text-white">
            Create Webhook Endpoint
          </DialogTitle>
        </DialogHeader>
        <div className="space-y-5 p-6">
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Name
            </Label>
            <Input
              placeholder="e.g. GitHub Deploy Hook"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="h-10 rounded-xl border-grey-200 bg-white text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Slug
            </Label>
            <div className="flex items-center gap-0 rounded-xl border border-grey-200 dark:border-white/10">
              <span className="bg-grey-50 px-3 py-2.5 text-xs text-grey-500 dark:bg-white/5 rounded-l-xl border-r border-grey-200 dark:border-white/10">
                /wh/
              </span>
              <Input
                placeholder="github-deploy"
                value={slug}
                onChange={(e) => setSlug(e.target.value)}
                className="h-10 border-0 bg-white text-sm focus-visible:ring-0 dark:bg-transparent rounded-l-none"
              />
            </div>
          </div>
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Target URL
            </Label>
            <Input
              placeholder="https://example.com/webhook"
              value={targetUrl}
              onChange={(e) => setTargetUrl(e.target.value)}
              className="h-10 rounded-xl border-grey-200 bg-white text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <div className="flex justify-end gap-2 pt-2">
            <Button
              variant="ghost"
              className="toss-press h-9 rounded-xl px-4 text-sm text-grey-600 dark:text-grey-400"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>
            <Button
              className="toss-press h-9 rounded-xl bg-blue-500 px-6 text-sm font-medium text-white hover:bg-blue-600"
              onClick={handleSubmit}
              disabled={createWebhook.isPending || !name || !slug}
            >
              {createWebhook.isPending ? 'Creating...' : 'Create Webhook'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
