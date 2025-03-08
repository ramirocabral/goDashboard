"use client"

import { useState } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { RefreshCw, ChevronDown, ChevronUp, HardDrive, Thermometer, AlertTriangle, CheckCircle } from "lucide-react"

const SmartData = () => {
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

  // Extract temperature from data if available
  const extractTemperature = (data) => {
    if (!data) return null

    // Look for temperature in various formats
    const tempKeys = Object.keys(data).filter(
      (key) => key.toLowerCase().includes("temperature") && !key.toLowerCase().includes("time"),
    )

    if (tempKeys.length > 0) {
      const tempValue = data[tempKeys[0]]
      // Extract numeric value if possible
      const match = tempValue.match(/(\d+)/)
      return match ? `${match[1]}°C` : tempValue
    }

    return null
  }

  // Extract health status or percentage used if available
  const extractHealthStatus = (data) => {
    if (!data) return { status: "Unknown", value: null }

    // Look for health status indicators
    if (data["Critical Warning"] && data["Critical Warning"].includes("0x00")) {
      return { status: "Healthy", value: "OK" }
    }

    if (data["Percentage Used"]) {
      const match = data["Percentage Used"].match(/(\d+)%/)
      const percentage = match ? Number.parseInt(match[1]) : null

      if (percentage !== null) {
        if (percentage < 50) {
          return { status: "Healthy", value: `${percentage}% Used` }
        } else if (percentage < 90) {
          return { status: "Warning", value: `${percentage}% Used` }
        } else {
          return { status: "Critical", value: `${percentage}% Used` }
        }
      }
    }

    // Look for other health indicators
    const healthKeys = ["Health Status", "SMART Health Status", "Available Spare", "Media and Data Integrity Errors"]

    for (const key of healthKeys) {
      if (data[key]) {
        const value = data[key].trim()
        if (value.includes("OK") || value.includes("PASSED") || value.includes("100%") || value === "0") {
          return { status: "Healthy", value }
        } else if (value.includes("WARN")) {
          return { status: "Warning", value }
        } else if (value.includes("FAIL") || value.includes("ERROR")) {
          return { status: "Critical", value }
        }
        return { status: "Unknown", value }
      }
    }

    return { status: "Unknown", value: null }
  }

  // Get status icon based on health status
  const getStatusIcon = (status) => {
    switch (status) {
      case "Healthy":
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case "Warning":
        return <AlertTriangle className="h-5 w-5 text-yellow-500" />
      case "Critical":
        return <AlertTriangle className="h-5 w-5 text-red-500" />
      default:
        return <HardDrive className="h-5 w-5 text-gray-400" />
    }
  }

  // If no data is available yet, show a loading state
  if (!smartData || !smartData.devices || smartData.devices.length === 0) {
    return (
      <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="rounded-lg bg-orange-500/10 p-2">
              <HardDrive className="h-5 w-5 text-orange-500" />
            </div>
            <div>
              <h3 className="text-sm font-medium text-gray-200">S.M.A.R.T.</h3>
              <p className="text-xs text-gray-400">Disk health information</p>
            </div>
          </div>
        </div>

        <div className="flex items-center justify-center h-40">
          <div className="text-center text-gray-400">
            <RefreshCw className="h-8 w-8 animate-spin mx-auto mb-2" />
            <p>Loading S.M.A.R.T. data...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-orange-500/10 p-2">
            <HardDrive className="h-5 w-5 text-orange-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">S.M.A.R.T.</h3>
            <p className="text-xs text-gray-400">{smartData.devices.length} disks monitored</p>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {smartData.devices.map((device, index) => {
          const temperature = extractTemperature(device.data)
          const health = extractHealthStatus(device.data)

          return (
            <div key={index} className="border border-gray-700 rounded-md overflow-hidden bg-gray-800/50">
              <div
                className="flex justify-between items-center p-3 cursor-pointer hover:bg-gray-700/50"
                onClick={() => toggleDisk(device.device)}
              >
                <div className="flex items-center">
                  {getStatusIcon(health.status)}
                  <div className="ml-2">
                    <div className="font-medium text-gray-200">{device.device}</div>
                    <div className="text-xs text-gray-400">{device.model || device.device.split("/").pop()}</div>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="text-right mr-4">
                    {health.value && (
                      <div className="font-medium text-gray-200">
                        {health.status}: {health.value}
                      </div>
                    )}
                    {temperature && (
                      <div className="text-xs text-gray-400 flex items-center justify-end">
                        <Thermometer className="h-3 w-3 mr-1" />
                        {temperature}
                      </div>
                    )}
                  </div>
                  {expandedDisk === device.device ? (
                    <ChevronUp className="h-5 w-5 text-gray-400" />
                  ) : (
                    <ChevronDown className="h-5 w-5 text-gray-400" />
                  )}
                </div>
              </div>

              {expandedDisk === device.device && device.data && (
                <div className="p-3 border-t border-gray-700 bg-gray-800/30">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                    {Object.entries(device.data).map(([key, value], i) => (
                      <div key={i} className="flex justify-between py-1 border-b border-gray-700/50">
                        <span className="text-xs text-gray-400">{key}</span>
                        <span className="text-xs text-gray-200 font-mono">{value}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default SmartData
