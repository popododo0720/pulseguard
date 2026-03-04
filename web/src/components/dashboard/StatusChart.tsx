import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import { useChartData } from '@/hooks/use-mock-data'

export function StatusChart() {
  const data = useChartData()

  return (
    <div className="animate-fade-in-up rounded-2xl bg-white p-6 dark:bg-[#161b22]" style={{ animationDelay: '200ms' }}>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h3 className="text-base font-semibold text-grey-900 dark:text-white">
            Execution Trend
          </h3>
          <p className="mt-0.5 text-sm text-grey-500">Last 7 days</p>
        </div>
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1.5">
            <div className="h-2 w-2 rounded-full bg-blue-500" />
            <span className="text-xs text-grey-500">Success</span>
          </div>
          <div className="flex items-center gap-1.5">
            <div className="h-2 w-2 rounded-full bg-red-500" />
            <span className="text-xs text-grey-500">Failure</span>
          </div>
        </div>
      </div>
      <div className="h-[240px]">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data} margin={{ top: 0, right: 0, left: -20, bottom: 0 }}>
            <defs>
              <linearGradient id="successGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3182f6" stopOpacity={0.15} />
                <stop offset="95%" stopColor="#3182f6" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="failureGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#f04452" stopOpacity={0.15} />
                <stop offset="95%" stopColor="#f04452" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#e5e8eb"
              vertical={false}
            />
            <XAxis
              dataKey="date"
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#8b95a1', fontSize: 12 }}
              dy={8}
            />
            <YAxis
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#8b95a1', fontSize: 12 }}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: '#191f28',
                border: 'none',
                borderRadius: '12px',
                padding: '10px 14px',
                boxShadow: 'none',
              }}
              itemStyle={{ color: '#e5e8eb', fontSize: '13px' }}
              labelStyle={{ color: '#8b95a1', fontSize: '12px', marginBottom: '4px' }}
              cursor={{ stroke: '#d1d6db', strokeDasharray: '4 4' }}
            />
            <Area
              type="monotone"
              dataKey="success"
              stroke="#3182f6"
              strokeWidth={2}
              fill="url(#successGradient)"
              dot={false}
              activeDot={{ r: 4, stroke: '#3182f6', strokeWidth: 2, fill: '#fff' }}
            />
            <Area
              type="monotone"
              dataKey="failure"
              stroke="#f04452"
              strokeWidth={2}
              fill="url(#failureGradient)"
              dot={false}
              activeDot={{ r: 4, stroke: '#f04452', strokeWidth: 2, fill: '#fff' }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}
