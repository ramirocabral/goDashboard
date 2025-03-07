"use client"

import { useState } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { RefreshCw, ChevronDown, ChevronUp, AlertTriangle, CheckCircle, AlertCircle } from "lucide-react"

const SmartWidget = () => {
  const { smartData, refreshStaticData } = useWebSocket()
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [expandedDisk, setExpandedDisk] = useState(null)

  const handleRefresh = async () => {
    setIsRefreshing(true)
    await refreshStaticData()
    setTimeout(() => setIsRefreshing(false), 500)
  }

  // Toggle expanded disk
  const toggleDisk = (diskName) => {
    if (expandedDisk === diskName) {
      setExpandedDisk(null)
    } else {
      setExpandedDisk(diskName)
    }
  }

  // Get status icon based on health status
  const getStatusIcon = (status) => {
    if (status === "PASSED" || status === "OK" || status === "GOOD") {
      return <CheckCircle className="h-5 w-5 text-green-500" />
    } else if (status === "WARNING") {
      return <AlertCircle className="h-5 w-5 text-yellow-500" />
    } else if (status === "FAILED" || status === "BAD") {
      return <AlertTriangle className="h-5 w-5 text-red-500" />
    } else {
      return <AlertCircle className="h-5 w-5 text-muted-foreground" />
    }
  }

  // If no data is available yet, show a loading state
  if (!smartData || !smartData.disks || smartData.disks.length === 0) {
    return (
      <div className="card">
        <div className="card-header">
          <div className="flex justify-between items-center">
            <div className="card-title">S.M.A.R.T. Data</div>
            <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
              <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
            </button>
          </div>
        </div>
        <div className="card-content">
          <div className="text-center text-muted-foreground py-4">Loading S.M.A.R.T. data...</div>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      <div className="card-header">
        <div className="flex justify-between items-center">
          <div className="card-title">S.M.A.R.T. Data</div>
          <button className="p-1 rounded-md hover:bg-secondary" onClick={handleRefresh}>
            <RefreshCw className={`h-4 w-4 ${isRefreshing ? "animate-spin" : ""}`} />
          </button>
        </div>
      </div>
      <div className="card-content">
        <div className="space-y-4">
          {smartData.disks.map((disk, index) => (
            <div key={index} className="border border-border rounded-md overflow-hidden">
              <div
                className="flex justify-between items-center p-3 cursor-pointer hover:bg-secondary/50"
                onClick={() => toggleDisk(disk.device)}
              >
                <div className="flex items-center">
                  {getStatusIcon(disk.health_status)}
                  <div className="ml-2">
                    <div className="font-medium">{disk.device}</div>
                    <div className="text-sm text-muted-foreground">{disk.model || "Unknown Model"}</div>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="text-right mr-4">
                    <div className="font-medium">Health: {disk.health_status}</div>
                    <div className="text-sm text-muted-foreground">Temperature: {disk.temperature}°C</div>
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
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div>
                      <div className="text-sm text-muted-foreground">Serial Number</div>
                      <div className="font-medium">{disk.serial_number || "N/A"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Firmware</div>
                      <div className="font-medium">{disk.firmware || "N/A"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Power On Hours</div>
                      <div className="font-medium">{disk.power_on_hours || "N/A"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Power Cycle Count</div>
                      <div className="font-medium">{disk.power_cycle_count || "N/A"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Interface</div>
                      <div className="font-medium">{disk.interface || "N/A"}</div>
                    </div>
                    <div>
                      <div className="text-sm text-muted-foreground">Rotation Rate</div>
                      <div className="font-medium">{disk.rotation_rate || "N/A"}</div>
                    </div>
                  </div>

                  {disk.attributes && disk.attributes.length > 0 && (
                    <div className="mt-4">
                      <div className="text-sm font-medium mb-2">S.M.A.R.T. Attributes</div>
                      <div className="overflow-x-auto">
                        <table className="w-full">
                          <thead>
                            <tr className="border-b border-border">
                              <th className="text-left pb-2">ID</th>
                              <th className="text-left pb-2">Name</th>
                              <th className="text-left pb-2">Value</th>
                              <th className="text-left pb-2">Worst</th>
                              <th className="text-left pb-2">Threshold</th>
                              <th className="text-left pb-2">Status</th>
                            </tr>
                          </thead>
                          <tbody>
                            {disk.attributes.map((attr, attrIndex) => (
                              <tr key={attrIndex} className="border-b border-border">
                                <td className="py-2">{attr.id}</td>
                                <td className="py-2">{attr.name}</td>
                                <td className="py-2">{attr.value}</td>
                                <td className="py-2">{attr.worst}</td>
                                <td className="py-2">{attr.threshold}</td>
                                <td className="py-2 flex items-center">
                                  {getStatusIcon(attr.status)}
                                  <span className="ml-1">{attr.status}</span>
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>
                    </div>
                  )}

                  {disk.self_test && (
                    <div className="mt-4">
                      <div className="text-sm font-medium mb-2">Self-Test Results</div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <div className="text-sm text-muted-foreground">Last Test Type</div>
                          <div className="font-medium">{disk.self_test.last_test_type || "N/A"}</div>
                        </div>
                        <div>
                          <div className="text-sm text-muted-foreground">Last Test Status</div>
                          <div className="font-medium flex items-center">
                            {getStatusIcon(disk.self_test.last_test_status)}
                            <span className="ml-1">{disk.self_test.last_test_status || "N/A"}</span>
                          </div>
                        </div>
                        <div>
                          <div className="text-sm text-muted-foreground">Last Test Time</div>
                          <div className="font-medium">{disk.self_test.last_test_time || "N/A"}</div>
                        </div>
                        <div>
                          <div className="text-sm text-muted-foreground">Lifetime Test Count</div>
                          <div className="font-medium">{disk.self_test.lifetime_test_count || "N/A"}</div>
                        </div>
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

export default SmartWidget

