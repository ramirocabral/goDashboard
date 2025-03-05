"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../../contexts/WebSocketContext"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { Network } from "lucide-react"

const NetworkCard = ({ splitView }) => {
  const { networkData } = useWebSocket()
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

  // Update realtime data when new network data arrives
  useEffect(() => {
    if (networkData) {
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            download: networkData.current?.rx_rate || 0,
            upload: networkData.current?.tx_rate || 0,
          },
        ]

        // Keep last 50 points
        if (newData.length > 50) {
          return newData.slice(-50)
        }
        return newData
      })
    }
  }, [networkData])

  return (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      {/* Icon and Title */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-yellow-500/10 p-2">
            <Network className="h-5 w-5 text-yellow-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Network</h3>
            <p className="text-xs text-gray-400">{networkData?.current?.interface || "eth0"}</p>
          </div>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-gray-200">{formatBytes(networkData?.current?.rx_rate || 0)}/s</p>
          <p className="text-xs text-gray-400">Current Speed</p>
        </div>
      </div>

      {/* Network Stats */}
      <div className="mb-4 grid grid-cols-2 gap-4">
        <div>
          <p className="text-xs text-gray-400">Download</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(networkData?.current?.rx_rate || 0)}/s</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Upload</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(networkData?.current?.tx_rate || 0)}/s</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Total Download</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(networkData?.current?.rx_total || 0)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Total Upload</p>
          <p className="text-sm font-medium text-gray-200">{formatBytes(networkData?.current?.tx_total || 0)}</p>
        </div>
      </div>

      {/* Graphs */}
      {splitView ? (
        <div className="grid grid-cols-2 gap-4">
          <div className="h-24">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={realtimeData}>
                <defs>
                  <linearGradient id="downloadGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="rgb(234, 179, 8)" stopOpacity={0.3} />
                    <stop offset="100%" stopColor="rgb(234, 179, 8)" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <Area
                  type="monotone"
                  dataKey="download"
                  stroke="rgb(234, 179, 8)"
                  strokeWidth={2}
                  fill="url(#downloadGradient)"
                  isAnimationActive={false}
                  dot={false}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
          <div className="h-24">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={realtimeData}>
                <defs>
                  <linearGradient id="uploadGradient" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="rgb(234, 179, 8)" stopOpacity={0.3} />
                    <stop offset="100%" stopColor="rgb(234, 179, 8)" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <Area
                  type="monotone"
                  dataKey="upload"
                  stroke="rgb(234, 179, 8)"
                  strokeWidth={2}
                  fill="url(#uploadGradient)"
                  isAnimationActive={false}
                  dot={false}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>
      ) : (
        <div className="h-24">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={realtimeData}>
              <defs>
                <linearGradient id="networkGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="rgb(234, 179, 8)" stopOpacity={0.3} />
                  <stop offset="100%" stopColor="rgb(234, 179, 8)" stopOpacity={0} />
                </linearGradient>
              </defs>
              <Area
                type="monotone"
                dataKey="download"
                stroke="rgb(234, 179, 8)"
                strokeWidth={2}
                fill="url(#networkGradient)"
                isAnimationActive={false}
                dot={false}
              />
              <Area
                type="monotone"
                dataKey="upload"
                stroke="rgb(234, 179, 8)"
                strokeWidth={1}
                strokeDasharray="3 3"
                fill="none"
                isAnimationActive={false}
                dot={false}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  )
}

export default NetworkCard

