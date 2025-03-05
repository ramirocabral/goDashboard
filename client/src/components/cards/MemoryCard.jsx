"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../../contexts/WebSocketContext"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { MemoryStickIcon as Memory } from "lucide-react"

const MemoryCard = () => {
  const { memoryData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState([])

  // Format bytes to human-readable format
  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B"
    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ["B", "KB", "MB", "GB", "TB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  // Update realtime data when new memory data arrives
  useEffect(() => {
    if (memoryData) {
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            value: (memoryData.current?.used / memoryData.current?.total) * 100 || 0,
            swap: (memoryData.current?.swap_used / memoryData.current?.swap_total) * 100 || 0,
          },
        ]

        // Keep last 50 points
        if (newData.length > 50) {
          return newData.slice(-50)
        }
        return newData
      })
    }
  }, [memoryData])

  return (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      {/* Icon and Title */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-red-500/10 p-2">
            <Memory className="h-5 w-5 text-red-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Memory</h3>
            <p className="text-xs text-gray-400">DDR4 3200MHz</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-gray-200">
            {((memoryData?.current?.used / memoryData?.current?.total) * 100)?.toFixed(1) || "0.0"}%
          </p>
          <p className="text-xs text-gray-400">
            {formatBytes(memoryData?.current?.used)} / {formatBytes(memoryData?.current?.total)}
          </p>
        </div>
      </div>

      {/* Memory Stats */}
      <div className="mb-4 grid grid-cols-2 gap-4">
        <div>
          <p className="text-xs text-gray-400">Used</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(memoryData?.current?.used)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Available</p>
          <p className="text-sm font-medium text-gray-200">
            {formatBytes(memoryData?.current?.total - memoryData?.current?.used)}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Swap Used</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(memoryData?.current?.swap_used)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Swap Total</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(memoryData?.current?.swap_total)}</p>
        </div>
      </div>

      {/* Graph */}
      <div className="h-24">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={realtimeData}>
            <defs>
              <linearGradient id="memoryGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="rgb(239, 68, 68)" stopOpacity={0.3} />
                <stop offset="100%" stopColor="rgb(239, 68, 68)" stopOpacity={0} />
              </linearGradient>
            </defs>
            <Area
              type="monotone"
              dataKey="value"
              stroke="rgb(239, 68, 68)"
              strokeWidth={2}
              fill="url(#memoryGradient)"
              isAnimationActive={false}
              dot={false}
            />
            <Area
              type="monotone"
              dataKey="swap"
              stroke="rgb(239, 68, 68, 0.5)"
              strokeWidth={1}
              strokeDasharray="3 3"
              fill="none"
              isAnimationActive={false}
              dot={false}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

export default MemoryCard

