"use client"

import { useState } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { RefreshCw, ChevronDown, ChevronUp, HardDrive } from "lucide-react"

const DisksInfo = () => {
  const { disksInfo, refreshStaticData } = useWebSocket()
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [expandedDisk, setExpandedDisk] = useState(null)

  const handleRefresh = async () => {
    setIsRefreshing(true)
    await refreshStaticData()
    setTimeout(() => setIsRefreshing(false), 500)
  }

  // Format bytes to human-readable format
  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B"

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ["B", "KB", "MB", "GB", "TB"]

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  // Calculate percentage
  const calculatePercentage = (used, total) => {
    if (!used || !total) return 0
    return Math.round((used / total) * 100)
  }

  // Toggle expanded disk
  const toggleDisk = (diskName) => {
    if (expandedDisk === diskName) {
      setExpandedDisk(null)
    } else {
      setExpandedDisk(diskName)
    }
  }

  // If no data is available yet, show a loading state
  if (!disksInfo || !disksInfo.disks || disksInfo.disks.length === 0) {
    return (
      <div className="card">
        <div className="card-header">
          <div className="flex justify-between items-center">
            <div className="card-title">Disks Information</div>
            <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
              <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
            </button>
          </div>
        </div>
        <div className="card-content">
          <div className="text-center text-muted-foreground py-4">Loading disks information...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      <div className="card-header">
        <div className="flex justify-between items-center">
          <div className="card-title">Disks Information</div>
          <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
            <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
          </button>
        </div>
      </div>
      <div className="card-content">
        <div className="space-y-4">
          {disksInfo.disks.map((disk, index) => (
            <div key={index} className="border border-border rounded-md overflow-hidden">
              <div
                className="flex justify-between items-center p-3 cursor-pointer hover:bg-secondary/50"
                onClick={() => toggleDisk(disk.device)}
              >
                <div className="flex items-center">
                  <HardDrive className="h-5 w-5 mr-2" />
                  <div>
                    <div className="font-medium">{disk.device}</div>
                    <div className="text-sm text-muted-foreground">{disk.model || "Unknown Model"}</div>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="text-right mr-4">
                    <div className="font-medium">{formatBytes(disk.size)}</div>
                    <div className="text-sm text-muted-foreground">
                      {calculatePercentage(disk.used, disk.size)}% Used
                    </div>
                  </div>
                  {expandedDisk === disk.device ? (
                    <ChevronUp className="h-5 w-5" />
                  ) : (
                    <ChevronDown className="h-5 w-5" />
                  )}
                </div>
              </div>

              {expandedDisk === disk.device && (
                <div className="p-3 border-t border-border bg-secondary/20">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                    <div>
                      <div className="text-sm text-muted-foreground">Mount Point</div>
                      <div className="font-medium">{disk.mount_point || "Not Mounted"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">File System</div>
                      <div className="font-medium">{disk.filesystem || "Unknown"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Used Space</div>
                      <div className="font-medium">
                        {formatBytes(disk.used)} ({calculatePercentage(disk.used, disk.size)}%)
                      </div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Free Space</div>
                      <div className="font-medium">{formatBytes(disk.size - disk.used)}</div>
                    </div>
                  </div>

                  <div className="w-full bg-secondary rounded-full h-2.5 mb-4">
                    <div
                      className="bg-primary h-2.5 rounded-full"
                      style={{ width: `${calculatePercentage(disk.used, disk.size)}%` }}
                    ></div>
                  </div>

                  {disk.partitions && disk.partitions.length > 0 && (
                    <div className="mt-4">
                      <div className="text-sm font-medium mb-2">Partitions</div>
                      <div className="overflow-x-auto">
                        <table className="w-full">
                          <thead>
                            <tr className="border-b border-border">
                              <th className="text-left pb-2">Name</th>
                              <th className="text-left pb-2">Mount Point</th>
                              <th className="text-left pb-2">Size</th>
                              <th className="text-left pb-2">Used</th>
                            </tr>
                          </thead>
                          <tbody>
                            {disk.partitions.map((partition, partIndex) => (
                              <tr key={partIndex} className="border-b border-border">
                                <td className="py-2">{partition.name}</td>
                                <td className="py-2">{partition.mount_point || "Not Mounted"}</td>
                                <td className="py-2">{formatBytes(partition.size)}</td>
                                <td className="py-2">
                                  {formatBytes(partition.used)} ({calculatePercentage(partition.used, partition.size)}%)
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default DisksInfo

