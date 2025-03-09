"use client"

import { useState, useEffect, memo } from "react"
import { format } from "date-fns"
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Legend, ResponsiveContainer, Tooltip } from "recharts"

const IoHistoryChart = ({ data }) => {
  const [chartData, setChartData] = useState([])
  const [devices, setDevices] = useState([])
  const [selectedDevice, setSelectedDevice] = useState("")
  const [yDomain, setYDomain] = useState([0, 1000])

  useEffect(() => {
    if (data) {
      const deviceList = data.devices.map((device) => device.device)
      setDevices(deviceList)

      if (!selectedDevice && deviceList.length > 0) {
        setSelectedDevice(deviceList[0])
      }
    }
  }, [data])

  useEffect(() => {
    if (selectedDevice) {
      const deviceData = data.devices.find((device) => device.device === selectedDevice)

      if (deviceData) {
        const newChartData = deviceData.stats.map((point) => ({
          time: new Date(point.timestamp),
          formattedTime: format(new Date(point.timestamp), "HH:mm:ss"),
          read: point.kb_read_per_second,
          write: point.kb_write_per_second,
        }))

        setChartData(newChartData)

        const values = newChartData.flatMap((d) => [d.read, d.write])
        const min = Math.min(...values, 0)
        const max = Math.round(Math.max(...values, 1) * 1.05)

        setYDomain([min, max])
      }
    }
  }, [selectedDevice, data])

  const formatXAxis = (timestamp) => format(new Date(timestamp), "HH:mm")

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="p-3 bg-gray-900 text-white rounded shadow-md">
          <p className="text-sm font-semibold">{format(new Date(label), "HH:mm:ss")}</p>
          <p className="text-sm text-blue-400">Read: {payload[0].value.toFixed(2)} KB/s</p>
          <p className="text-sm text-green-400">Write: {payload[1].value.toFixed(2)} KB/s</p>
        </div>
      )
    }
    return null
  }

  return (
    <div className="bg-widget p-6 rounded-xl shadow-md overflow-x-auto">
    
      <div className="mb-4">
        <label className="text-white text-sm">Select Device:</label>
        <select
          className="ml-2 p-2 bg-gray-800 text-white rounded-md"
          value={selectedDevice}
          onChange={(e) => setSelectedDevice(e.target.value)}
        >
          {devices.map((device) => (
            <option key={device} value={device}>{device}</option>
          ))}
        </select>
      </div>

      <div className="h-[400px]">
        <ResponsiveContainer width="100%" height="100%" minWidth={500}>
          <LineChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#444" />
            <XAxis tick={{ dy: 7}} dataKey="time" tickFormatter={formatXAxis} stroke="#ddd" />
            <YAxis domain={yDomain} stroke="#ddd" />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line type="monotone" dataKey="read" name="Read (KB/s)" stroke="#4F46E5" strokeWidth={2} />
            <Line type="monotone" dataKey="write" name="Write (KB/s)" stroke="#22C55E" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

export default IoHistoryChart
