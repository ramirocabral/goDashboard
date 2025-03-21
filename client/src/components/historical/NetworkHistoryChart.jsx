"use client";

import { useState, useEffect } from "react";
import { format } from "date-fns";
import {
  LineChart,
  Line,
  XAxis,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
} from "recharts";

const NetworkHistoryChart = ({ data }) => {
  const [chartData, setChartData] = useState([]);
  const [interfaces, setInterfaces] = useState([]);
  const [selectedInterface, setSelectedInterface] = useState("");

  useEffect(() => {
    if (data) {
      const interfaceList = data.interfaces.map((iface) => iface.interface);
      setInterfaces(interfaceList);

      if (!selectedInterface && interfaceList.length > 0) {
        setSelectedInterface(interfaceList[0]);
      }
    }
  }, [data]);

  useEffect(() => {
    if (selectedInterface) {
      if (!selectedInterface || !data) return;

      const interfaceData = data.interfaces.find(
        (iface) => iface.interface === selectedInterface,
      );
      if (!interfaceData || !interfaceData.stats) return;

      if (interfaceData) {
        const newChartData = interfaceData.stats.map((point) => ({
          time: new Date(point.timestamp),
          formattedTime: format(new Date(point.timestamp), "HH:mm:ss"),
          rxBytes: point.rx_bytes_ps,
          txBytes: point.tx_bytes_ps,
        }));

        setChartData(newChartData);
      }
    }
  }, [selectedInterface, data]);

  const formatBytes = (bytes, decimals = 2) => {
    if (!bytes) return "0 B";
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["B", "KB", "MB", "GB", "TB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${Number.parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
  };

  const formatXAxis = (timestamp) => format(new Date(timestamp), "HH:mm");

  const CustomTooltip = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
      return (
        <div className="p-3 bg-gray-900 text-white rounded shadow-md">
          <p className="text-sm font-semibold">
            {format(new Date(label), "HH:mm:ss")}
          </p>
          <p className="text-sm text-blue-400">
            Rx: {formatBytes(payload[0].value)}/s
          </p>
          <p className="text-sm text-green-400">
            Tx: {formatBytes(payload[1].value)}/s
          </p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="bg-widget p-6 rounded-xl shadow-md overflow-x-auto">
      <div className="mb-4">
        <label className="text-white text-sm">Select Interface:</label>
        <select
          className="ml-2 p-2 bg-gray-800 text-white rounded-md"
          value={selectedInterface}
          onChange={(e) => setSelectedInterface(e.target.value)}
        >
          {interfaces.map((iface) => (
            <option key={iface} value={iface}>
              {iface}
            </option>
          ))}
        </select>
      </div>

      <div className="h-[400px]">
        <ResponsiveContainer width="100%" height="100%" minWidth={500}>
          <LineChart
            data={chartData}
            margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="#444" />
            <XAxis
              tick={{ dy: 7 }}
              dataKey="time"
              tickFormatter={formatXAxis}
              stroke="#ddd"
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line
              type="monotone"
              dataKey="rxBytes"
              name="Rx (B/s)"
              stroke="#4F46E5"
              strokeWidth={2}
            />
            <Line
              type="monotone"
              dataKey="txBytes"
              name="Tx (B/s)"
              stroke="#22C55E"
              strokeWidth={2}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default NetworkHistoryChart;
