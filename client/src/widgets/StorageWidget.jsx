"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { HardDrive, ChevronDown, ChevronUp } from "lucide-react"
import CardContainer from "../components/widgets/WidgetContainer"
import WidgetGrid from "../components/widgets/WidgetGrid"

const StorageCard = () => {
  const { disksInfo } = useWebSocket()
  const [selectedDisk, setSelectedDisk] = useState(null)
  const [showDiskSelector, setShowDiskSelector] = useState(false)

  const disks = disksInfo?.disks || [
  ]

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

  return (
    <CardContainer>
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-emerald-500/10 p-2">
            <HardDrive className="h-5 w-5 text-emerald-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Storage</h3>
            <div className="relative">
              <button
                onClick={() => setShowDiskSelector(!showDiskSelector)}
                className="flex items-center text-xs text-gray-400 hover:text-gray-300"
              >
                {selectedDisk?.device || "Select Disk"}
                {disks.length > 1 &&
                  (showDiskSelector ? (
                    <ChevronUp className="ml-1 h-3 w-3" />
                  ) : (
                    <ChevronDown className="ml-1 h-3 w-3" />
                  ))}
              </button>

              {showDiskSelector && disks.length > 1 && (
                <div className="absolute top-full left-0 z-10 mt-1 w-40 rounded-md bg-gray-800 shadow-lg">
                  <ul className="py-1">
                    {disks.map((disk) => (
                      <li key={disk.device}>
                        <button
                          className={`block w-full px-4 py-2 text-left text-xs ${
                            disk.device === selectedDisk?.device
                              ? "bg-gray-700 text-gray-200"
                              : "text-gray-400 hover:bg-gray-700 hover:text-gray-200"
                          }`}
                          onClick={() => {
                            setSelectedDisk(disk)
                            setShowDiskSelector(false)
                          }}
                        >
                          {disk.device}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        </div>
        {selectedDisk && (
          <div className="text-right">
            <p className="text-2xl font-bold text-gray-200">{selectedDisk.used_percentage}%</p>
            <p className="text-xs text-gray-400">{formatGB(selectedDisk.gb_used)} / {formatGB(selectedDisk.gb_size)}</p>
          </div>
        )}
      </div>

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
              {/* Progress circle */}
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

export default StorageCard
