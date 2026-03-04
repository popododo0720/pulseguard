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
import { useCreateJob, useAgents } from '@/hooks/use-api'

export function CreateJobDialog() {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')
  const [schedule, setSchedule] = useState('')
  const [command, setCommand] = useState('')
  const [agentId, setAgentId] = useState('')
  const createJob = useCreateJob()
  const { data: agents } = useAgents()

  const handleSubmit = () => {
    if (!name || !schedule || !command) return
    createJob.mutate(
      { name, schedule, command, agent_id: agentId || undefined } as never,
      {
        onSuccess: () => {
          setOpen(false)
          setName('')
          setSchedule('')
          setCommand('')
          setAgentId('')
        },
      },
    )
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="toss-press h-9 rounded-xl bg-blue-500 px-4 text-sm font-medium text-white hover:bg-blue-600">
          <Plus className="mr-1.5 h-4 w-4" />
          New Job
        </Button>
      </DialogTrigger>
      <DialogContent className="rounded-2xl border-grey-200 bg-white p-0 sm:max-w-lg dark:border-white/10 dark:bg-[#161b22]">
        <DialogHeader className="border-b border-grey-100 p-6 dark:border-white/5">
          <DialogTitle className="text-lg font-semibold text-grey-900 dark:text-white">
            Create New Job
          </DialogTitle>
        </DialogHeader>
        <div className="space-y-5 p-6">
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Job Name
            </Label>
            <Input
              placeholder="e.g. Database Backup"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="h-10 rounded-xl border-grey-200 bg-white text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Schedule (cron expression)
            </Label>
            <Input
              placeholder="e.g. 0 2 * * *"
              value={schedule}
              onChange={(e) => setSchedule(e.target.value)}
              className="h-10 rounded-xl border-grey-200 bg-white font-mono text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Command
            </Label>
            <Input
              placeholder="e.g. /opt/scripts/backup.sh"
              value={command}
              onChange={(e) => setCommand(e.target.value)}
              className="h-10 rounded-xl border-grey-200 bg-white font-mono text-sm dark:border-white/10 dark:bg-white/5"
            />
          </div>
          <div className="space-y-2">
            <Label className="text-sm font-medium text-grey-700 dark:text-grey-300">
              Agent
            </Label>
            <select
              value={agentId}
              onChange={(e) => setAgentId(e.target.value)}
              className="h-10 w-full rounded-xl border border-grey-200 bg-white px-3 text-sm dark:border-white/10 dark:bg-white/5 dark:text-grey-300"
            >
              <option value="">Select an agent...</option>
              {agents?.map((agent) => (
                <option key={agent.id} value={agent.id}>
                  {agent.hostname || agent.name}
                </option>
              ))}
            </select>
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
              disabled={createJob.isPending || !name || !schedule || !command}
            >
              {createJob.isPending ? 'Creating...' : 'Create Job'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
