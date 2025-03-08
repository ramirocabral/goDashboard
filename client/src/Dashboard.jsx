"use client"

import CpuWidget from "./widgets/CpuWidget"
import MemoryWidget from "./widgets/MemoryWidget"
import StorageWidget from "./widgets/StorageWidget"
import NetworkWidget from "./widgets/NetworkWidget"
import Containers from "./widgets/ContainersWidget"
import SmartWidget from "./widgets/SmartWidget"
import { useWebSocket } from "./contexts/WebSocketContext"
import SystemWidget from "./widgets/SystemWidget"
import IoWidget from "./widgets/IoWidget"

const Dashboard = () => {
  const { connected } = useWebSocket()

  const isConnected = Object.values(connected).some((status) => status)

  return (
    <div className="container mx-auto p-4">
      <header className="flex justify-between items-center mb-6">
        <div className="flex items-center">
          <h1 className="text-3xl font-bold">goDashboard</h1>
          <div
            className={`ml-4 flex items-center ${
              isConnected ? "text-green-500" : "text-red-500"
            }`}
          >
            <div
              className={`h-2 w-2 rounded-full ${
                isConnected ? "bg-green-500 animate-pulse" : "bg-red-500"
              }`}
            />
            <span className="ml-2 text-sm">
              {isConnected ? "Connected" : "Disconnected"}
            </span>
          </div>
        </div>
        <div className="text-sm text-white-600">
          {new Date().toLocaleDateString(undefined, {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
          })}
        </div>
      </header>

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
  );
}

export default Dashboard

