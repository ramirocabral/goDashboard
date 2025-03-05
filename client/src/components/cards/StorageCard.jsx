"use client"

import { useState, useEffect } from "react"
import { useWebSocket } from "../../contexts/WebSocketContext"
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts"
import TimeRangeSelector from "../TimeRangeSelector"
import { subSeconds } from "date-fns"

const StorageCard = () => {
  const { ioData, fetchHistoricalData } = useWebSocket()
  const [historicalData, setHistoricalData] = useState([])
  const [selectedRange, setSelectedRange] = useState(3600) // 1 hour by default
  const [showHistorical, setShowHistorical] = useState(false)

  // Placeholder data for storage info
  const storageInfo = {
    drives: [
      { name: "INTEL NVMe", size: 953.9, unit: "GiB" },
      { name: "SanDisk HD", size: 59.6, unit: "GiB" },
    ],
  }

  // Placeholder data for pie chart
  const pieData = [
    { name: "Used", value: 70 },
    { name: "Free", value: 30 },
  ]

  const COLORS = ["hsl(var(--disk-color))", "#555"]

  // Fetch historical data when range changes
  useEffect(() => {
    if (showHistorical) {
      const fetchData = async () => {
        const endTime = new Date().toISOString()
        const startTime = subSeconds(new Date(), selectedRange).toISOString()
        const interval = Math.max(Math.floor(selectedRange / 60), 1) // At least 1 second interval

        const data = await fetchHistoricalData("io", startTime, endTime, interval)
        if (data && data.devices && data.devices.length > 0) {
          setHistoricalData(
            data.devices[0].stats.map((item) => ({
              time: new Date(item.timestamp).getTime(),
              read: item.kb_read_per_second,
              write: item.kb_write_per_second,
            })),
          )
        }
      }

      fetchData()
    }
  }, [selectedRange, showHistorical, fetchHistoricalData])

  // Custom tooltip for pie chart
  const CustomTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-card p-2 border border-border rounded shadow-sm">
          <p className="text-sm">{`${payload[0].name}: ${payload[0].value}%`}</p>
        </div>
      )
    }
    return null
  }

  return (
    <div className="card">
      <div className="card-header">
        <div className="flex justify-between items-center">
          <div className="flex items-center">
            <div className="w-8 h-8 rounded-full bg-[hsl(var(--disk-color))] flex items-center justify-center mr-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-5 w-5 text-white"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M21 5H3a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h18a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2Z"></path>
                <path d="M18 12H6"></path>
              </svg>
            </div>
            <div className="card-title">Storage</div>
          </div>
          <button className="text-sm text-primary hover:underline" onClick={() => setShowHistorical(!showHistorical)}>
            {showHistorical ? "Hide I/O Stats" : "Show I/O Stats"}
          </button>
        </div>
      </div>
      <div className="card-content">
        <div className="grid grid-cols-1 gap-4 mb-4">
          {storageInfo.drives.map((drive, index) => (
            <div key={index} className="flex justify-between items-center">
              <div>
                <div className="text-sm text-muted-foreground">Drive</div>
                <div className="font-medium">{drive.name}</div>
              </div>
              <div className="text-right">
                <div className="text-sm text-muted-foreground">&nbsp;</div>
                <div className="font-medium">
                  {drive.size} {drive.unit}
                </div>
              </div>
            </div>
          ))}
        </div>

        {showHistorical ? (
          <>
            <div className="mb-4">
              <TimeRangeSelector selectedRange={selectedRange} onRangeChange={setSelectedRange} />
            </div>
            <div className="mt-4">
              <h3 className="text-sm font-medium mb-2">Disk I/O</h3>
              <div className="space-y-2">
                {historicalData.length > 0 ? (
                  <>
                    <div className="flex justify-between">
                      <span className="text-sm">Read:</span>
                      <span className="text-sm font-medium">
                        {historicalData[historicalData.length - 1]?.read || 0} KB/s
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-sm">Write:</span>
                      <span className="text-sm font-medium">
                        {historicalData[historicalData.length - 1]?.write || 0} KB/s
                      </span>
                    </div>
                  </>
                ) : (
                  <div className="text-sm text-muted-foreground">No I/O data available</div>
                )}
              </div>
            </div>
          </>
        ) : (
          <div className="graph-container flex items-center justify-center">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={40}
                  outerRadius={60}
                  paddingAngle={5}
                  dataKey="value"
                >
                  {pieData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip content={<CustomTooltip />} />
              </PieChart>
            </ResponsiveContainer>
          </div>
        )}

        {!showHistorical && (
          <div className="mt-2 flex justify-between items-center">
            <div className="stat-value">70%</div>
            <div className="stat-label">Disk Usage</div>
          </div>
        )}
      </div>
    </div>
  )
}

export default StorageCard

