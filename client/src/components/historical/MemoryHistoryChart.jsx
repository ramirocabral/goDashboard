"use client"

import { useState, useEffect, memo } from "react"
import { format } from "date-fns"
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Legend, ResponsiveContainer, Tooltip } from "recharts"

const MemoryHistoryChart = ({ data }) => {
  const [chartData, setChartData] = useState([])
  const [memoryInfo, setMemoryInfo] = useState({})

  useEffect(() => {
    if (data) {
      setMemoryInfo({
        type: data.type,
        frequency: data.frequency,
      })

      // console.log(data.stats)

      setChartData(
        data.stats.map((point) => ({
          time: new Date(point.timestamp),
          formattedTime: format(new Date(point.timestamp), "HH:mm:ss"),
          usage: point.used_percentage,
        }))
      )
    }
  }, [data])

  const formatXAxis = (timestamp) => format(new Date(timestamp), "HH:mm")

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="p-3 bg-gray-900 text-white rounded shadow-md">
          <p className="text-sm font-semibold">{format(new Date(label), "HH:mm:ss")}</p>
          <p className="text-sm text-blue-400">Memory Usage: {payload[0].value.toFixed(2)}%</p>
        </div>
      )
    }
    return null
  }

  return (
    <div className="bg-widget p-6 rounded-xl shadow-md overflow-x-auto">
      <div className="flex justify-between sm:grid sm:grid-cols-2 gap-4 mb-6 rounded-xl md:px-5 w-full">
        <InfoCard label="Type" value={memoryInfo.type || "N/A"} />
        <InfoCard label="Frequency" value={memoryInfo.frequency ? `${memoryInfo.frequency} MHz` : "N/A"} />
      </div>

      <div className="h-[400px]">
        <ResponsiveContainer width="100%" height="100%" minWidth={500}>
          <LineChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#444" />
            <XAxis tick={{ dy: 7}} dataKey="time" tickFormatter={formatXAxis} stroke="#ddd" />
            <YAxis domain={[0, 100]} stroke="#ddd" />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line type="monotone" dataKey="usage" name="Memory Usage (%)" stroke="#4F46E5" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}

const InfoCard = ({ label, value }) => (
  <div className="p-4  bg-gray-900 text-white rounded shadow w-full">
    <div className="text-sm text-gray-400">{label}</div>
    <div className="text-lg font-semibold">{value}</div>
  </div>
)

export default MemoryHistoryChart
