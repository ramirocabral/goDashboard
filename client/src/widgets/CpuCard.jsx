"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { Cpu } from "lucide-react"
import CardContainer from "../components/cards/CardContainer"
import CardHeader from "../components/cards/CardHeader"
import Chart from "../components/cards/Chart"
import InfoGrid from "../components/cards/InfoGrid"

const CpuWidget = () => {
  const { cpuData } = useWebSocket()
  const [realtimeData, setRealtimeData] = useState([])

  useEffect(() => {
    if (cpuData) {
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            value: cpuData.usage?.usage_percentage || 0,
          },
        ]
        //keep last 50 points
        if (newData.length > 50) {
          return newData.slice(-50)
        }
        return newData
      })
    }
  }, [cpuData])

  return (
    <CardContainer>
      <CardHeader 
        icon={<Cpu className="h-5 w-5 text-blue-500" />}
        title="Processor"
        subtitle={cpuData?.model_name || "NULL"} 
        value={(cpuData?.usage?.usage_percentage?.toFixed(1) || "0.0") + "%"}
        secondValue={(cpuData?.temp || "NULL") + "°C"}
      />
      <InfoGrid
        data={[
          { label: "Cores", value: cpuData?.cores || "NULL" },
          { label: "Threads", value: cpuData?.threads || "NULL" },
          { label: "Frequency", value: cpuData?.info?.frequency || "2.6" + " GHz" },
          { label: "Architecture", value: "x64" },
        ]}
       />
      <Chart 
        realTimeData={realtimeData} 
        id="cpuGradient"
        stopColor="rgb(59, 130, 246)"
        dataKey="value"
        stroke="rgb(59, 130, 246)"
        strokeWidth={2}
        fill="url(#cpuGradient)"
      />
    </CardContainer>
  )
}

export default CpuWidget

