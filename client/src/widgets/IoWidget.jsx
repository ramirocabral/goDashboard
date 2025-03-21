"use client";

import { useState, useEffect } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";
import { HardDrive } from "lucide-react";
import { Area, AreaChart, ResponsiveContainer } from "recharts";
import WidgetContainer from "../components/widgets/WidgetContainer";
import WidgetGrid from "../components/widgets/WidgetGrid";
import WidgetHeaderSelector from "../components/widgets/WidgetHeaderSelector";

const IoWidget = () => {
  const { ioData } = useWebSocket();
  const [realtimeData, setRealtimeData] = useState({});
  const [selectedDevice, setSelectedDevice] = useState(null);

  const devices = ioData || [];

  useEffect(() => {
    if (devices.length > 0 && !selectedDevice) {
      setSelectedDevice(devices[0].device);
    }
  }, [devices, selectedDevice]);

  const currentDeviceData = devices.find(
    (device) => device.device === selectedDevice,
  ) ||
    devices[0] || {
      device: "N/A",
      kb_read_per_second: 0,
      kb_write_per_second: 0,
    };

  useEffect(() => {
    if (ioData && devices.length > 0) {
      setRealtimeData((prevData) => {
        const newData = { ...prevData };

        devices.forEach((device) => {
          const deviceName = device.device;
          const readRate = device.kb_read_per_second || 0;
          const writeRate = device.kb_write_per_second || 0;

          if (!newData[deviceName]) {
            newData[deviceName] = [];
          }

          const timestamp = Date.now();
          newData[deviceName] = [
            ...newData[deviceName],
            {
              time: timestamp,
              read: readRate,
              write: writeRate,
            },
          ];

          if (newData[deviceName].length > 30) {
            newData[deviceName] = newData[deviceName].slice(-30);
          }
        });

        return newData;
      });
    }
  }, [ioData, devices]);

  const selectedDeviceData = realtimeData[selectedDevice] || [];

  const formatKB = (kb) => {
    if (kb < 1024) {
      return `${kb} KB/s`;
    } else if (kb < 1024 * 1024) {
      return `${(kb / 1024).toFixed(2)} MB/s`;
    } else {
      return `${(kb / (1024 * 1024)).toFixed(2)} GB/s`;
    }
  };

  const ioInfoData = [
    {
      label: "Read",
      value: formatKB(currentDeviceData.kb_read_per_second || 0),
    },
    {
      label: "Write",
      value: formatKB(currentDeviceData.kb_write_per_second || 0),
    },
  ];

  return (
    <WidgetContainer>
      <WidgetHeaderSelector
        icon={
          <div className="rounded-lg bg-purple-500/10 p-2">
            <HardDrive className="h-5 w-5 text-purple-500" />
          </div>
        }
        title="Disk I/O"
        items={devices.map((device) => device.device)}
        selectedItem={currentDeviceData.device}
        onItemSelect={setSelectedDevice}
        valueText={formatKB(currentDeviceData.kb_write_per_second || 0)}
        valueSubtext="Write Speed"
      />

      <WidgetGrid data={ioInfoData} />

      <div className="h-32">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={selectedDeviceData}>
            <defs>
              <linearGradient id="readGradient" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="0%"
                  stopColor="rgb(168, 85, 247)"
                  stopOpacity={0.3}
                />
                <stop
                  offset="100%"
                  stopColor="rgb(168, 85, 247)"
                  stopOpacity={0}
                />
              </linearGradient>
              <linearGradient id="writeGradient" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="0%"
                  stopColor="rgb(139, 92, 246)"
                  stopOpacity={0.3}
                />
                <stop
                  offset="100%"
                  stopColor="rgb(139, 92, 246)"
                  stopOpacity={0}
                />
              </linearGradient>
            </defs>
            <Area
              type="monotone"
              dataKey="read"
              stroke="rgb(168, 85, 247)"
              strokeWidth={2}
              fill="url(#readGradient)"
              isAnimationActive={false}
              dot={false}
            />
            <Area
              type="monotone"
              dataKey="write"
              stroke="rgb(139, 92, 246)"
              strokeWidth={2}
              fill="url(#writeGradient)"
              isAnimationActive={false}
              dot={false}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      <div className="mt-2 flex justify-between items-center">
        <div className="flex items-center">
          <div className="w-3 h-0.5 bg-purple-500 mr-1"></div>
          <span className="text-xs text-gray-400">Read</span>
        </div>
        <div className="flex items-center">
          <div className="w-3 h-0.5 bg-violet-500 mr-1"></div>
          <span className="text-xs text-gray-400">Write</span>
        </div>
      </div>
    </WidgetContainer>
  );
};

export default IoWidget;
