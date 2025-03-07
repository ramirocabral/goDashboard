"use client"

import React from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { RefreshCw } from "lucide-react"

const HostInfo = () => {
  const { uptimeData } = useWebSocket()
  const { hostInfo, refreshStaticData } = useWebSocket()
  const [isRefreshing, setIsRefreshing] = React.useState(false)

  // Format uptime
  const formatUptime = (seconds) => {
    if (!seconds) return "0 minutes"

    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)

    let result = ""
    if (days > 0) result += `${days} day${days > 1 ? "s" : ""} `
    if (hours > 0) result += `${hours} hour${hours > 1 ? "s" : ""} `
    if (minutes > 0) result += `${minutes} minute${minutes > 1 ? "s" : ""}`

    return result.trim()
  }


  const handleRefresh = async () => {
    setIsRefreshing(true)
    await refreshStaticData()
    setTimeout(() => setIsRefreshing(false), 500)
  }

  // If no data is available yet, show a loading state
  if (!hostInfo) {
    return (
      <div className="card">
        <div className="card-header">
          <div className="flex justify-between items-center">
            <div className="card-title">Host Information</div>
            <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
              <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
            </button>
          </div>
        </div>
        <div className="card-content">
          <div className="text-center text-muted-foreground py-4">Loading host information...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      <div className="card-header">
        <div className="flex justify-between items-center">
          <div className="card-title">Host Information</div>
          <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
            <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
          </button>
        </div>
      </div>
      <div className="card-content">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <div className="text-sm text-muted-foreground">Hostname</div>
            <div className="font-medium">{hostInfo.hostname || "N/A"}</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">OS</div>
            <div className="font-medium">{hostInfo.os || "N/A"}</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Kernel</div>
            <div className="font-medium">{hostInfo.kernel || "N/A"}</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Uptime</div>
            <div className="font-medium">{formatUptime(uptimeData?.uptime)}</div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default HostInfo