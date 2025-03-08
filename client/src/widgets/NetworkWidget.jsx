"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { Network, ChevronDown, ChevronUp } from "lucide-react"
import CardContainer from "../components/widgets/WidgetContainer"
import Chart from "../components/widgets/WidgetChart"
import WidgetGrid from "../components/widgets/WidgetGrid"
import WidgetHeaderSelector from "../components/widgets/WidgetHeaderSelector"

const NetworkWidget = () => {
  const { networkData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState({})
  const [selectedInterface, setSelectedInterface] = useState(null)

  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B"
    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ["B", "KB", "MB", "GB", "TB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
  }

  const interfaces = networkData || []

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

          if (!newData[interfaceName]) {
            newData[interfaceName] = []
          }

          const timestamp = Date.now()
          newData[interfaceName] = [
            ...newData[interfaceName],
            {
              time: timestamp,
              download: downloadRate,
              upload: uploadRate,
            },
          ]

          if (newData[interfaceName].length > 30) {
            newData[interfaceName] = newData[interfaceName].slice(-30)
          }
        })

        return newData
      })
    }
  }, [networkData, interfaces])

  const selectedInterfaceData = realtimeData[selectedInterface] || []

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

  if (interfaces.length > 1) {
    networkInfoData.push(
      { label: "Interfaces", value: interfaces.length.toString() },
      { label: "Active", value: interfaces.length.toString() },
    )
  }

  return (
    <CardContainer>

      <WidgetHeaderSelector
        icon={
        <div className="rounded-lg bg-yellow-500/10 p-2">
          <Network className="h-5 w-5 text-yellow-500" />
        </div>
        }
        title="Network"
        items={interfaces.map((iface) => iface.interface)}
        selectedItem={currentInterfaceData.interface}
        onItemSelect={setSelectedInterface}
        valueText={formatBytes(currentInterfaceData.usage?.rx_bytes_ps || 0) + "/s"}
        valueSubtext="Current Download"
      />

      <WidgetGrid data={networkInfoData} />

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