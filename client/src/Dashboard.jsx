"use client"

import CpuWidget from "./widgets/CpuCard"
import MemoryWidget from "./widgets/MemoryWidget"
import StorageWidget from "./widgets/StorageWidget"
import NetworkCard from "./widgets/NetworkCard"
import Containers from "./widgets/ContainersWidget"
import SmartWidget from "./widgets/SmartWidget"
import { useWebSocket } from "./contexts/WebSocketContext"
import SystemWidget from "./widgets/SystemWidget"

const Dashboard = () => {
  const { connected } = useWebSocket()

  const isConnected = Object.values(connected).some((status) => status)

  return (
    <div className="container mx-auto p-4">
      <header className="flex justify-between items-center mb-6">
        <div className="flex items-center">
          <h1 className="text-3xl font-bold">goMonitor</h1>
          <div className={`ml-4 flex items-center ${isConnected ? "text-green-500" : "text-red-500"}`}>
            <div className={`h-2 w-2 rounded-full ${isConnected ? "bg-green-500 animate-pulse" : "bg-red-500"}`} />
            <span className="ml-2 text-sm">{isConnected ? "Connected" : "Disconnected"}</span>
          </div>
        </div>
        <div className="text-sm text-white-600">
          {new Date().toLocaleDateString(undefined, { weekday: "long", year: "numeric", month: "long", day: "numeric" })}
        </div>
      </header>

      <div className="grid-container mb-6 grid-cols-2">
        <div className="align-self-start">
          <SystemWidget />
        </div>
        {/* <div className="col-span-2"> */}
          <Containers/>
        {/* </div> */}
      </div>

      <div className="grid-container mb-6">
        <CpuWidget/>
        <MemoryWidget />
        <NetworkCard/>
        <StorageWidget />
      </div>

      <div className="mb-6">
        <SmartWidget />
      </div>
    </div>
  )
}

export default Dashboard

