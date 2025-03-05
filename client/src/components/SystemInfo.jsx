import { useWebSocket } from "../contexts/WebSocketContext"

const SystemInfo = () => {
  const { cpuData, uptimeData } = useWebSocket()

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

  return (
    <div className="card">
      <div className="card-header">
        <div className="card-title">System Information</div>
      </div>
      <div className="card-content">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <div className="text-sm text-muted-foreground">OS</div>
            <div className="font-medium">Ubuntu 22.04 LTS</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Architecture</div>
            <div className="font-medium">x64</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Uptime</div>
            <div className="font-medium">{formatUptime(uptimeData?.uptime)}</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">CPU</div>
            <div className="font-medium">{cpuData?.model_name || "Loading..."}</div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">Cores / Threads</div>
            <div className="font-medium">
              {cpuData?.cores || "0"} / {cpuData?.threads || "0"}
            </div>
          </div>
          <div>
            <div className="text-sm text-muted-foreground">CPU Temperature</div>
            <div className="font-medium">{cpuData?.temp ? `${cpuData.temp}°C` : "N/A"}</div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default SystemInfo

