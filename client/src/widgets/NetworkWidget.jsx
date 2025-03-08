"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { Network, ChevronDown, ChevronUp } from "lucide-react"
import CardContainer from "../components/widgets/WidgetContainer"
import Chart from "../components/widgets/WidgetChart"
import InfoGrid from "../components/widgets/WidgetGrid"

const NetworkWidget = () => {
  const { networkData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState({})
  const [selectedInterface, setSelectedInterface] = useState(null)
  const [showInterfaceSelector, setShowInterfaceSelector] = useState(false)

  // Format bytes to human-readable format
  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B"
    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ["B", "KB", "MB", "GB", "TB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  // Get available interfaces
  const interfaces = networkData || []

  // Set default selected interface if not set
  useEffect(() => {
    if (interfaces.length > 0 && !selectedInterface) {
      setSelectedInterface(interfaces[0].interface)
    }
  }, [interfaces, selectedInterface])

  const currentInterfaceData = interfaces.find((iface) => iface.interface === selectedInterface) ||
    interfaces[0] || { interface: "N/A", usage: { rx_bytes_ps: 0, tx_bytes_ps: 0 } }

  useEffect(() => {
    if (networkData && interfaces.length > 0) {

      setRealtimeData((prevData) => {
        const newData = { ...prevData }

        interfaces.forEach((iface) => {
          const interfaceName = iface.interface
          const downloadRate = iface.usage?.rx_bytes_ps || 0
          const uploadRate = iface.usage?.tx_bytes_ps || 0

          // Initialize array if it doesn't exist
          if (!newData[interfaceName]) {
            newData[interfaceName] = []
          }

          // Add new data point
          const timestamp = Date.now()
          newData[interfaceName] = [
            ...newData[interfaceName],
            {
              time: timestamp,
              download: downloadRate,
              upload: uploadRate,
            },
          ]

          // Keep only the last 30 points to make the graph smoother
          if (newData[interfaceName].length > 30) {
            newData[interfaceName] = newData[interfaceName].slice(-30)
          }
        })

        return newData
      })
    }
  }, [networkData, interfaces])

  // Get data for the selected interface
  const selectedInterfaceData = realtimeData[selectedInterface] || []

  // Prepare data for InfoGrid component
  const networkInfoData = [
    {
      label: "Download",
      value: formatBytes(currentInterfaceData.usage?.rx_bytes_ps || 0) + "/s",
    },
    {
      label: "Upload",
      value: formatBytes(currentInterfaceData.usage?.tx_bytes_ps || 0) + "/s",
    },
  ]

  // Add interface count info if there are multiple interfaces
  if (interfaces.length > 1) {
    networkInfoData.push(
      { label: "Interfaces", value: interfaces.length.toString() },
      { label: "Active", value: interfaces.length.toString() },
    )
  }

  return (
    <CardContainer>
      {/* Custom header to avoid nesting issues */}
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-yellow-500/10 p-2">
            <Network className="h-5 w-5 text-yellow-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Network</h3>
            {/* Interface selector as a sibling, not a child of p */}
            <div className="relative">
              <button
                onClick={() => setShowInterfaceSelector(!showInterfaceSelector)}
                className="flex items-center text-xs text-gray-400 hover:text-gray-300"
              >
                {currentInterfaceData?.interface}
                {interfaces.length > 1 &&
                  (showInterfaceSelector ? (
                    <ChevronUp className="ml-1 h-3 w-3" />
                  ) : (
                    <ChevronDown className="ml-1 h-3 w-3" />
                  ))}
              </button>

              {/* Interface selector dropdown */}
              {showInterfaceSelector && interfaces.length > 1 && (
                <div className="absolute top-full left-0 z-10 mt-1 w-40 rounded-md bg-gray-800 shadow-lg">
                  <ul className="py-1">
                    {interfaces.map((iface) => (
                      <li key={iface.interface}>
                        <button
                          className={`block w-full px-4 py-2 text-left text-xs ${
                            iface.interface === selectedInterface
                              ? "bg-gray-700 text-gray-200"
                              : "text-gray-400 hover:bg-gray-700 hover:text-gray-200"
                          }`}
                          onClick={() => {
                            setSelectedInterface(iface.interface)
                            setShowInterfaceSelector(false)
                          }}
                        >
                          {iface.interface}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-gray-200">
            {formatBytes(currentInterfaceData.usage?.rx_bytes_ps || 0)}/s
          </p>
          <p className="text-xs text-gray-400">Current Download</p>
        </div>
      </div>

      {/* Network Stats */}
      <InfoGrid data={networkInfoData} />

      {/* Single chart for network data */}
      <div className="h-32">
        <Chart
          realTimeData={selectedInterfaceData}
          id="networkGradient"
          stopColor="rgb(234, 179, 8)"
          dataKey="download"
          stroke="rgb(234, 179, 8)"
          strokeWidth={2}
          fill="url(#networkGradient)"
        />
      </div>
    </CardContainer>
  )
}

export default NetworkWidget