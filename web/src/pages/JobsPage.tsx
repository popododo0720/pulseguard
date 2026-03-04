import { useState } from 'react'
import { JobList } from '@/components/jobs/JobList'
import { JobDetail } from '@/components/jobs/JobDetail'
import { CreateJobDialog } from '@/components/jobs/CreateJobDialog'
import type { Job } from '@/hooks/use-api'

export function JobsPage() {
  const [selectedJob, setSelectedJob] = useState<Job | null>(null)

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight text-grey-900 dark:text-white">
            Jobs
          </h1>
          <p className="mt-1 text-sm text-grey-500">
            Manage and monitor your scheduled cron jobs
          </p>
        </div>
        <CreateJobDialog />
      </div>

      {/* Content */}
      <div className={selectedJob ? 'grid gap-6 lg:grid-cols-[1fr_400px]' : ''}>
        <JobList
          onSelectJob={setSelectedJob}
          selectedJobId={selectedJob?.id}
        />
        {selectedJob && (
          <div className="hidden lg:block">
            <div className="sticky top-24 max-h-[calc(100vh-8rem)] overflow-hidden">
              <JobDetail
                job={selectedJob}
                onClose={() => setSelectedJob(null)}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
