"use client"

import { useState } from "react"
import { useWebSocket } from "../contexts/WebSocketContext"
import { RefreshCw, ChevronDown, ChevronUp, CheckCircle, Logs } from "lucide-react"

const SmartData = () => {
  const { smartData } = useWebSocket()
  const [expandedDisk, setExpandedDisk] = useState(null)

  const toggleDisk = (diskName) => {
    if (expandedDisk === diskName) {
      setExpandedDisk(null)
    } else {
      setExpandedDisk(diskName)
    }
  }

  const formatSataData = (data) => {
    if (!data || !data["ID#"]) return ""

    const headerLine = data["ID#"]
    const columns = headerLine.split(/\s+/).filter(Boolean)

    const columnWidths = columns.map((col, index) => {
      let maxWidth = col.length
      Object.entries(data).forEach(([key, value]) => {
        if (key === "ID#") return
        const values = value.split(/\s+/).filter(Boolean)
        if (values[index] && values[index].length > maxWidth) {
          maxWidth = values[index].length
        }
      })
      return maxWidth
    })

    const formattedHeader = columns.map((col, i) => col.padEnd(columnWidths[i])).join(" ")

    const rows = Object.entries(data)
      .filter(([key]) => key !== "ID#")
      .sort((a, b) => Number.parseInt(a[0]) - Number.parseInt(b[0]))
      .map(([_, value]) => {
        const values = value.split(/\s+/).filter(Boolean)
        return values.map((val, i) => val.padEnd(i < columnWidths.length ? columnWidths[i] : 0)).join(" ")
      })

    return [formattedHeader, ...rows].join("\n")
  }

const formatNvmeData = (data) => {
  if (!data) return "";

  // find the longest key to insert spaces
  const maxKeyLength = Math.max(...Object.keys(data).map(k => k.length));

  return Object.entries(data)
    .map(([key, value]) => {
      const cleanedValue = value.trim();
      return `${key.padEnd(maxKeyLength, " ")} : ${cleanedValue}`;
    })
    .join("\n");
};

  // loading state
  if (!smartData || !smartData.devices || smartData.devices.length === 0) {
    return (
      <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="rounded-lg bg-orange-500/10 p-2">
              <Logs className=" h-5 w-5 text-orange-500"/>
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
            <Logs className=" h-5 w-5 text-orange-500"/>
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">S.M.A.R.T.</h3>
            <p className="text-xs text-gray-400">{smartData.devices.length} disks monitored</p>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {smartData.devices.map((device, index) => {
          const deviceName = device.device.split("/").pop()

          return (
            <div key={index} className="border border-gray-700 rounded-md overflow-hidden bg-gray-800/50">
              <div
                className="flex justify-between items-center p-3 cursor-pointer hover:bg-gray-700/50"
                onClick={() => toggleDisk(device.device)}
              >
                <div className="flex items-center">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <div className="ml-2">
                    <div className="font-medium text-gray-200">{device.device}</div>
                    <div className="text-xs text-gray-400">{device.type === "sata" ? "SATA Drive" : "NVMe Drive"}</div>
                  </div>
                </div>
                <div className="flex items-center">
                  <div className="text-right mr-4">
                  </div>
                  {expandedDisk === device.device ? (
                    <ChevronUp className="h-5 w-5 text-gray-400" />
                  ) : (
                    <ChevronDown className="h-5 w-5 text-gray-400" />
                  )}
                </div>
              </div>

              {expandedDisk === device.device && (
                <div className="p-3 border-t border-gray-700 bg-gray-800/30">
                <pre className="font-mono text-xs text-gray-200 whitespace-pre overflow-x-auto bg-[#1a1b26] p-4 rounded">
                  {device.type === "sata" ? formatSataData(device.data) : formatNvmeData(device.data)}
                </pre>
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

