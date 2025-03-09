"use client";

import { useState, useEffect } from "react";
import { format } from "date-fns";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

const CpuHistoryChart = ({ data }) => {
  const [chartData, setChartData] = useState([]);
  const [cpuInfo, setCpuInfo] = useState({});

  useEffect(() => {
    if (data) {
      setCpuInfo({
        model: data.model_name,
        frequency: data.frequency,
        cores: data.cores,
        threads: data.threads,
      });

      setChartData(
        data.stats.map((point) => ({
          time: new Date(point.timestamp),
          formattedTime: format(new Date(point.timestamp), "HH:mm:ss"),
          usage: point.usage_percentage,
          idle: point.idle_percentage,
        }))
      );
    }
  }, [data]);

  const formatXAxis = (timestamp) => format(new Date(timestamp), "HH:mm");

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="p-3 bg-gray-900 text-white rounded shadow-md">
          <p className="text-sm font-semibold">
            {format(new Date(label), "HH:mm:ss")}
          </p>
          <p className="text-sm text-blue-400">
            CPU Usage: {payload[0].value.toFixed(2)}%
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="bg-widget p-6 rounded-xl shadow-md overflow-x-auto">
      <div className="flex justify-between md:grid md:grid-cols-4 gap-4 mb-6 rounded-xl md:px-5">
        <InfoCard label="CPU Model" value={cpuInfo.model || "N/A"} />
        <InfoCard
          label="Frequency"
          value={cpuInfo.frequency ? `${cpuInfo.frequency} MHz` : "N/A"}
        />
        <InfoCard label="Cores" value={cpuInfo.cores || "N/A"} />
        <InfoCard label="Threads" value={cpuInfo.threads || "N/A"} />
      </div>

      <div className="h-[400px]">
        <ResponsiveContainer width="100%" height="100%" minWidth={500}>
          <LineChart
            data={chartData}
            margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="#444" />
            <XAxis tick={{ dy: 7}} dataKey="time" tickFormatter={formatXAxis} stroke="#ddd" />
            <YAxis domain={[0, 100]} stroke="#ddd" />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line
              type="monotone"
              dataKey="usage"
              name="CPU Usage (%)"
              stroke="#4F46E5"
              strokeWidth={2}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

const InfoCard = ({ label, value }) => (
  <div className="p-4  bg-gray-900 text-white rounded shadow">
    <div className="text-sm text-gray-400">{label}</div>
    <div className="text-lg font-semibold">{value}</div>
  </div>
);

export default CpuHistoryChart;
