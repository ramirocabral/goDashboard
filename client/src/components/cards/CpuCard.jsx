"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../../contexts/WebSocketContext"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { Cpu } from "lucide-react"

const CpuCard = ({ showAllCores }) => {
  const { cpuData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState([])

  // Update realtime data when new CPU data arrives
  useEffect(() => {
    if (cpuData) {
      // console.log("nashee")
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            value: cpuData.usage?.usage_percentage || 0,
          },
        ]
        // Keep last 50 points
        if (newData.length > 50) {
          return newData.slice(-50)
        }
        return newData
      })
    }
  }, [cpuData])

  return (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      {/* Icon and Title */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-blue-500/10 p-2">
            <Cpu className="h-5 w-5 text-blue-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Processor</h3>
            <p className="text-xs text-gray-400">{cpuData?.model_name || "NULL"}</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-gray-200">{cpuData?.usage?.usage_percentage?.toFixed(1) || "0.0"}%</p>
          <p className="text-xs text-gray-400">{cpuData?.temp || "NULL"}°C</p>
        </div>
      </div>

      {/* Specs Grid */}
      <div className="mb-4 grid grid-cols-2 gap-4">
        <div>
          <p className="text-xs text-gray-400">Cores</p>
          <p className="text-sm font-medium text-gray-200">{cpuData?.cores || "NULL"}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Threads</p>
          <p className="text-sm font-medium text-gray-200">{cpuData?.threads || "NULL"}</p>
        </div>
        {/* <div>
          <p className="text-xs text-gray-400">Frequency</p>
          <p className="text-sm font-medium text-gray-200">{cpuData?.info?.frequency || "2.6"} GHz</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Architecture</p>
          <p className="text-sm font-medium text-gray-200">x64</p>
        </div> */}
      </div>

      {/* Graph */}
      <div className="h-24">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={realtimeData}>
            <defs>
              <linearGradient id="cpuGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="rgb(59, 130, 246)" stopOpacity={0.3} />
                <stop offset="100%" stopColor="rgb(59, 130, 246)" stopOpacity={0} />
              </linearGradient>
            </defs>
            <Area
              type="monotone"
              dataKey="value"
              stroke="rgb(59, 130, 246)"
              strokeWidth={2}
              fill="url(#cpuGradient)"
              isAnimationActive={false}
              dot={false}
            />
            {showAllCores &&
              cpuData?.info?.cores > 0 &&
              Array.from({ length: cpuData.info.cores }).map((_, i) => (
                <Area
                  key={i}
                  type="monotone"
                  dataKey={`core${i}`}
                  stroke={`rgb(59, 130, 246, ${0.3 + i * 0.1})`}
                  strokeWidth={1}
                  fill="none"
                  isAnimationActive={false}
                  dot={false}
                />
              ))}
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

export default CpuCard

