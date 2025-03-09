"use client";

import CpuWidget from "../widgets/CpuWidget";
import MemoryWidget from "../widgets/MemoryWidget";
import StorageWidget from "../widgets/StorageWidget";
import NetworkWidget from "../widgets/NetworkWidget";
import Containers from "../widgets/ContainersWidget";
import SmartWidget from "../widgets/SmartWidget";
import { useWebSocket } from "../contexts/WebSocketContext";
import SystemWidget from "../widgets/SystemWidget";
import IoWidget from "../widgets/IoWidget";
import Navbar from "../components/navbar/Navbar";

const Dashboard = () => {
  const { connected } = useWebSocket();

  const isConnected = Object.values(connected).some((status) => status);

  return (
    <div className="min-h-screen bg-background text-foreground transition-colors duration-300">
      <div className="container mx-auto p-4">

      <Navbar />

        <div className="grid grid-cols-1 lg:grid-cols-3 mb-6 gap-4">
          <div className="lg:col-span-1 h-full">
            <SystemWidget />
          </div>
          <div className="lg:col-span-2">
            <Containers />
          </div>
        </div>

        <div className="grid-container mb-6">
          <CpuWidget />
          <MemoryWidget />
          <NetworkWidget />
          <StorageWidget />
        </div>

        <div className="container">
          <div className="flex flex-col lg:flex-row gap-4 mb-6">
            <div className="lg:w-2/3">
              <SmartWidget />
            </div>
            <div className="lg:w-1/3">
              <IoWidget />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
