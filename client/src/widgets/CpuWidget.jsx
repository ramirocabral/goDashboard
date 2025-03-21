"use client";

import { useState, useEffect } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";
import { Cpu } from "lucide-react";
import CardContainer from "../components/widgets/WidgetContainer";
import CardHeader from "../components/widgets/WidgetHeader";
import Chart from "../components/widgets/WidgetChart";
import WidgetGrid from "../components/widgets/WidgetGrid";

const CpuWidget = () => {
  const { cpuData } = useWebSocket();
  const [realtimeData, setRealtimeData] = useState([]);

  useEffect(() => {
    if (cpuData) {
      setRealtimeData((prev) => {
        const newData = [
          ...prev,
          {
            time: Date.now(),
            value: cpuData.usage?.usage_percentage || 0,
          },
        ];
        if (newData.length > 50) {
          return newData.slice(-50);
        }
        return newData;
      });
    }
  }, [cpuData]);

  return (
    <CardContainer>
      <CardHeader
        icon={<Cpu className="h-5 w-5 text-blue-500" />}
        title="Processor"
        subtitle={cpuData?.model_name || "NULL"}
        value={(cpuData?.usage?.usage_percentage?.toFixed(1) || "0.0") + "%"}
        secondValue={(cpuData?.temp || "NULL") + "°C"}
      />
      <WidgetGrid
        data={[
          { label: "Frequency", value: cpuData?.frequency || "2.6" + " GHz" },
          { label: "Family", value: cpuData?.family || "x64" },
          { label: "Cores", value: cpuData?.cores || "NULL" },
          { label: "Threads", value: cpuData?.threads || "NULL" },
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
  );
};

export default CpuWidget;
