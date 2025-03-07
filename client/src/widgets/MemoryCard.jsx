"use client"

import { useState, useEffect } from "react"
import { Area, AreaChart, ResponsiveContainer } from "recharts"
import { MemoryStickIcon as Memory } from "lucide-react"
import { useWebSocket } from "../contexts/WebSocketContext"
import CardContainer from "../components/cards/CardContainer"
import CardHeader from "../components/cards/CardHeader"
import Chart from "../components/cards/Chart"
import InfoGrid from "../components/cards/InfoGrid"


const MemoryCard = () => {
  const { memoryData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState([])

  //data received is in kb
  const formatBytes = (kbytes, decimals = 2) => {
    if (!kbytes) return "0 B"
    const sizes = ["KB", "MB", "GB", "TB"]
    const i = Math.floor(Math.log(kbytes) / Math.log(1024))
    return `${parseFloat((kbytes / Math.pow(1024, i)).toFixed(decimals))} ${sizes
      [i]}
    `
  }

  // Update realtime data when new memory data arrives
  useEffect(() => {
    if (memoryData) {
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            value: (memoryData.used / memoryData.total) * 100 || 0,
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
    <CardContainer>
      <CardHeader
        icon={<Memory className="h-5 w-5 text-red-500" />}
        title="Memory"
        subtitle="DDR4 3200MHz"
        value={(((memoryData?.used / memoryData?.total) * 100)?.toFixed(1) || "0.0") + "%"}
        secondValue={`${formatBytes(memoryData?.used)} / ${formatBytes(memoryData?.total)}`} 
      />
      <InfoGrid
        data={[
          { label: "Active", value: formatBytes(memoryData?.active) },
          { label: "Inactive", value: formatBytes(memoryData?.inactive) },
          { label: "Buffers", value: formatBytes(memoryData?.buffers) },
          { label: "Cached", value: formatBytes(memoryData?.cached) },
        ]}
      />
      <Chart
        realTimeData={realtimeData}
        id="memoryGradient"
        stopColor="rgb(239, 68, 68)"
        dataKey="value"
        stroke="rgb(239, 68, 68)"
        strokeWidth={2}
        fill="url(#memoryGradient)"
      />
    </CardContainer>
  )
}

export default MemoryCard

