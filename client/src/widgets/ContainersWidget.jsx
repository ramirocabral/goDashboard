"use client"

import { useState } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { Box, Play, Pause, RotateCcw, Search, ChevronDown, ChevronUp } from "lucide-react"

const Containers = () => {
  const { containerData } = useWebSocket()
  const [filter, setFilter] = useState("")
  const [expanded, setExpanded] = useState(false)

  const filteredContainers =
    containerData?.filter(
      (container) =>
        container.name.toLowerCase().includes(filter.toLowerCase()) ||
        container.image.toLowerCase().includes(filter.toLowerCase()) ||
        container.id.toLowerCase().includes(filter.toLowerCase()),
    ) || []

  const formatUptime = (seconds) => {
    if (!seconds) return "N/A"

    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)

    if (days > 0) {
      return `${days}d ${hours}h`
    } else if (hours > 0) {
      return `${hours}h ${minutes}m`
    } else {
      return `${minutes}m`
    }
  }

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case "running":
        return "text-green-500"
      case "paused":
        return "text-yellow-500"
      case "stopped":
        return "text-red-500"
      default:
        return "text-muted-foreground"
    }
  }

  const getStatusBgColor = (status) => {
    switch (status.toLowerCase()) {
      case "running":
        return "bg-green-500/20"
      case "paused":
        return "bg-yellow-500/20"
      case "stopped":
        return "bg-red-500/20"
      default:
        return "bg-gray-500/20"
    }
  }

  return (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      {/* Icon and Title */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-cyan-500/10 p-2">
            <Box className="h-5 w-5 text-cyan-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Containers</h3>
            <p className="text-xs text-gray-400">{containerData?.containers?.length || 0} containers</p>
          </div>
        </div>
        <div className="flex items-center space-x-2">
          <div className="relative">
            <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Filter..."
              className="w-32 rounded-md border border-gray-700 bg-gray-800 py-1 pl-8 pr-2 text-xs text-gray-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
          <button className="rounded-md bg-gray-800 p-1 hover:bg-gray-700" onClick={() => setExpanded(!expanded)}>
            {expanded ? (
              <ChevronUp className="h-4 w-4 text-gray-400" />
            ) : (
              <ChevronDown className="h-4 w-4 text-gray-400" />
            )}
          </button>
        </div>
      </div>

      {!containerData ? (
        <div className="flex h-40 items-center justify-center">
          <div className="text-center text-gray-400">
            <Box className="mx-auto mb-2 h-8 w-8 animate-pulse" />
            <p>Loading container data...</p>
          </div>
        </div>
      ) : (
        <div className="overflow-hidden">
          {filteredContainers.length === 0 ? (
            <div className="py-4 text-center text-gray-400">No containers found</div>
          ) : (
            <div className="space-y-2">
              {filteredContainers.slice(0, expanded ? filteredContainers.length : 3).map((container) => (
                <div
                  key={container.id}
                  className="flex items-center justify-between rounded-md bg-gray-800/50 p-2 hover:bg-gray-800"
                >
                  <div className="flex items-center space-x-2">
                    <div
                      className={`h-2 w-2 rounded-full ${container.status.toLowerCase() === "running" ? "bg-green-500" : container.status.toLowerCase() === "paused" ? "bg-yellow-500" : "bg-red-500"}`}
                    ></div>
                    <div>
                      <p className="text-sm font-medium text-gray-200">{container.name}</p>
                      <p className="text-xs text-gray-400">{container.image}</p>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="text-right">
                      <p className="text-xs text-gray-400">CPU: {container.cpu_usage}%</p>
                      <p className="text-xs text-gray-400">MEM: {container.memory_usage} MB</p>
                    </div>
                    <div
                      className={`rounded-md px-2 py-1 text-xs ${getStatusBgColor(container.status)} ${getStatusColor(container.status)}`}
                    >
                      {container.status}
                    </div>
                    <div className="flex space-x-1">
                      {container.status.toLowerCase() !== "running" ? (
                        <button className="rounded-md p-1 text-green-500 hover:bg-gray-700">
                          <Play className="h-3 w-3" />
                        </button>
                      ) : (
                        <button className="rounded-md p-1 text-yellow-500 hover:bg-gray-700">
                          <Pause className="h-3 w-3" />
                        </button>
                      )}
                      <button className="rounded-md p-1 text-blue-500 hover:bg-gray-700">
                        <RotateCcw className="h-3 w-3" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}

              {filteredContainers.length > 3 && !expanded && (
                <button
                  className="mt-2 w-full rounded-md bg-gray-800 py-1 text-xs text-gray-400 hover:bg-gray-700"
                  onClick={() => setExpanded(true)}
                >
                  Show {filteredContainers.length - 3} more containers
                </button>
              )}

              {expanded && filteredContainers.length > 3 && (
                <button
                  className="mt-2 w-full rounded-md bg-gray-800 py-1 text-xs text-gray-400 hover:bg-gray-700"
                  onClick={() => setExpanded(false)}
                >
                  Show less
                </button>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default Containers