import { Search, Moon, Sun, Bell } from 'lucide-react'
import { useEffect, useState } from 'react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export function Header() {
  const [isDark, setIsDark] = useState(false)

  useEffect(() => {
    const root = document.documentElement
    if (isDark) {
      root.classList.add('dark')
    } else {
      root.classList.remove('dark')
    }
  }, [isDark])

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-grey-200 bg-grey-50/80 px-8 backdrop-blur-sm dark:border-white/[0.06] dark:bg-[#0d1117]/80">
      {/* Search */}
      <div className="relative w-full max-w-md">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-grey-400" />
        <Input
          placeholder="Search jobs, webhooks, agents..."
          className="h-9 border-grey-200 bg-white pl-10 text-sm text-grey-700 placeholder:text-grey-400 focus-visible:ring-blue-500/20 dark:border-white/10 dark:bg-white/5 dark:text-grey-200"
        />
      </div>

      {/* Actions */}
      <div className="flex items-center gap-1">
        <Button
          variant="ghost"
          size="icon"
          className="toss-press h-9 w-9 text-grey-500 hover:bg-grey-100 hover:text-grey-700 dark:hover:bg-white/5"
        >
          <Bell className="h-[18px] w-[18px]" strokeWidth={1.8} />
        </Button>
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
