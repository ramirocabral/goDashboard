import { useEffect, useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import CpuHistoryChart from "./CpuHistoryChart";
import MemoryHistoryChart from "./MemoryHistoryChart";
import IoHistoryChart from "./IoHistoryChart";
import NetworkHistoryChart from "./NetworkHistoryChart";
import TimeRangeSelector from "./TimeRangeSelector";
import { API_BASE_URL_HISTORY } from "../../constants/endpoints";

const HistoricalMetrics = () => {
  const [cpuData, setCpuData] = useState(null);
  const [memoryData, setMemoryData] = useState(null);
  const [ioData, setIoData] = useState(null);
  const [networkData, setNetworkData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState("cpu");
  const [timeRange, setTimeRange] = useState("1h");

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const end = new Date().toISOString();
        let start, interval;
        switch (timeRange) {
          case "1h":
            start = new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString();
            interval = "1m";
            break;
          case "3h":
            start = new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString();
            interval = "3m";
            break;
          case "6h":
            start = new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString();
            interval = "5m";
            break;
          case "12h":
            start = new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString();
            interval = "10m";
            break;
          default:
            start = new Date(Date.now() - 1 * 60 * 60 * 1000).toISOString();
            interval = "1m";
        }

        const urls = ["cpu", "memory", "io", "network"].map(
          (metric) => `${API_BASE_URL_HISTORY}/${metric}?start=${start}&end=${end}&interval=${interval}`
        );

        const [cpuResponse, memoryResponse, ioResponse, networkResponse] = await Promise.all(
          urls.map((url) => fetch(url).then((res) => res.json()))
        );

        setCpuData(cpuResponse);
        setMemoryData(memoryResponse);
        setIoData(ioResponse);
        setNetworkData(networkResponse);
      } catch (error) {
        console.error("Error fetching historical data:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [timeRange]);

  const TabButton = ({ label, tab }) => (
    <button
      className={`px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap ${
        activeTab === tab
          ? "border-primary text-primary"
          : "border-transparent text-muted-foreground hover:text-foreground hover:border-border"
      }`}
      onClick={() => setActiveTab(tab)}
    >
      {label}
    </button>
  );

  return (
    <Card className="px-2 py-4 sm:px-4 border-border">
      <div className="flex flex-col sm:flex-row sm:justify-between items-center mb-4">
        <div className="flex space-x-4 border-b pb-2 border-border overflow-x-auto w-full sm:w-auto">
          <TabButton label="CPU" tab="cpu" />
          <TabButton label="Memory" tab="memory" />
          <TabButton label="Disk I/O" tab="io" />
          <TabButton label="Network" tab="network" />
        </div>

      <TimeRangeSelector timeRange={timeRange} setTimeRange={setTimeRange} />
      </div> 

      <CardContent className="mt-4">
        {loading ? (
          <p>Loading...</p>
        ) : (
          <div className="w-full overflow-x-auto">
            {activeTab === "cpu" && cpuData && <CpuHistoryChart data={cpuData} />}
            {activeTab === "memory" && memoryData && <MemoryHistoryChart data={memoryData} />}
            {activeTab === "io" && ioData && <IoHistoryChart data={ioData} />}
            {activeTab === "network" && networkData && <NetworkHistoryChart data={networkData} />}
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default HistoricalMetrics;
