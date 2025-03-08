"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { HardDrive} from "lucide-react"
import CardContainer from "../components/widgets/WidgetContainer"
import WidgetGrid from "../components/widgets/WidgetGrid"
import WidgetHeaderSelector from "../components/widgets/WidgetHeaderSelector"

const StorageWidget = () => {
  const { disksInfo } = useWebSocket()
  const [selectedDisk, setSelectedDisk] = useState(null)

  const disks = disksInfo?.disks || []

  useEffect(() => {
    if (!selectedDisk && disks.length > 0) {
      setSelectedDisk(disks[0])
    }
  }, [disks, selectedDisk])

  const calculateStrokeDasharray = (percentage, radius) => {
    const circumference = 2 * Math.PI * radius
    const dashArray = (percentage / 100) * circumference
    return `${dashArray} ${circumference}`
  }

  const formatGB = (gb) => {
    if (gb >= 1000) {
      return `${(gb / 1000).toFixed(1)} TB`
    }
    return `${gb.toFixed(1)} GB`
  }

  //we need this shit because of the json response format
  const handleDiskSelect = (device) => {
    setSelectedDisk(disks.find((disk) => disk.device === device))
  }

  return (
    <CardContainer>

      <WidgetHeaderSelector
        icon={
          <div className="rounded-lg bg-emerald-500/10 p-2">
            <HardDrive className="h-5 w-5 text-emerald-500" />
          </div>
        }
        title="Storage"
        items={disks.map((disk) => disk.device)}
        selectedItem={selectedDisk?.device}
        onItemSelect={handleDiskSelect}
        valueText={selectedDisk? `${selectedDisk?.used_percentage}%`: "0%"}
        valueSubtext={selectedDisk? `${formatGB(selectedDisk?.gb_used)} / ${formatGB(selectedDisk?.gb_size)}` : "0 GB"}
      />

      {selectedDisk && (
        <div className="flex flex-col justify-center">
          <WidgetGrid data={[
            { label: "Type", value: selectedDisk.type },
            { label: "Mount Point", value: selectedDisk.mount_point },
            { label: "Used Space", value: formatGB(selectedDisk.gb_used) },
            { label: "Free Space", value: formatGB(selectedDisk.gb_free) }
          ]} />

          <div className="relative h-48 w-48 -mt-5 -mb-6 mx-auto">
            <svg className="h-full w-full -rotate-90 ">
              <circle
                cx="96"
                cy="96"
                r="50"
                stroke="currentColor"
                strokeWidth="15"
                fill="none"
                className="text-gray-700"
              />
              
              <circle
                cx="96"
                cy="96"
                r="50"
                stroke="currentColor"
                strokeWidth="15"
                fill="none"
                className="text-emerald-500"
                strokeLinecap="round"
                strokeDasharray={calculateStrokeDasharray(selectedDisk.used_percentage, 50)}
              />
            </svg>
            {/* Center text */}
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <span className="text-2xl font-bold text-gray-200">{selectedDisk.used_percentage}%</span>
              <span className="text-xs text-gray-400">Used</span>
            </div>
          </div>
       </div>
      )}
    </CardContainer>
  )
}

export default StorageWidget
