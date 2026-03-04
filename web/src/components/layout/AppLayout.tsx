import { Outlet } from 'react-router-dom'
import { Sidebar } from './Sidebar'
import { Header } from './Header'

export function AppLayout() {
  return (
    <div className="flex min-h-screen bg-grey-50 dark:bg-[#0d1117]">
      <Sidebar />
      <div className="flex flex-1 flex-col pl-[240px]">
        <Header />
        <main className="flex-1 px-8 py-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
