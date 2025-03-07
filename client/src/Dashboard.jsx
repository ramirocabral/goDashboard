"use client"

import { useState } from "react"
import CpuCard from "./widgets/CpuCard"
import MemoryCard from "./widgets/MemoryCard"
import StorageCard from "./widgets/StorageCard"
import NetworkCard from "./widgets/NetworkCard"
import Containers from "./widgets/Containers"
import ContainerList from "./components/ContainerList"
// import SystemInfo from "./SystemInfo"
import HostInfo from "./components/HostInfo"
import DisksInfo from "./components/DisksInfo"
import SmartData from "./components/SmartData"
import { useWebSocket } from "./contexts/WebSocketContext"
import SysInfo from "./components/SysInfo"

const Dashboard = ({ darkMode, setDarkMode }) => {
  const { connected } = useWebSocket()
  const [splitView, setSplitView] = useState(false)
  const [showAllCores, setShowAllCores] = useState(false)

  // Calculate overall connection status
  const isConnected = Object.values(connected).some((status) => status)

  return (
    <div className="container mx-auto p-4">
      <header className="flex justify-between items-center mb-6">
        <div className="flex items-center">
          <h1 className="text-3xl font-bold">dash.</h1>
          <div className={`ml-4 flex items-center ${isConnected ? "text-green-500" : "text-red-500"}`}>
            <div className={`h-2 w-2 rounded-full ${isConnected ? "bg-green-500 animate-pulse" : "bg-red-500"}`} />
            <span className="ml-2 text-sm">{isConnected ? "Connected" : "Disconnected"}</span>
          </div>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center">
            <span className="mr-2 text-sm">Dark Mode</span>
            <label className="switch">
              <input
                type="checkbox"
                className="switch-input"
                checked={darkMode}
                onChange={() => setDarkMode(!darkMode)}
              />
              <span className="switch-slider"></span>
            </label>
          </div>
          <div className="flex items-center">
            <span className="mr-2 text-sm">Split View</span>
            <label className="switch">
              <input
                type="checkbox"
                className="switch-input"
                checked={splitView}
                onChange={() => setSplitView(!splitView)}
              />
              <span className="switch-slider"></span>
            </label>
          </div>
        </div>
      </header>


      {/* <div className="mb-6">
        <SystemInfo />
      </div> */}

      <div className="grid-container mb-6 grid-cols-3">
        <SysInfo />
        <HostInfo />
        <Containers />
      </div>

      <div className="grid-container mb-6">
        <CpuCard showAllCores={showAllCores} />
        <MemoryCard />
        <StorageCard />
        <NetworkCard splitView={splitView} />
      </div>

      <div className="mb-6">
        <DisksInfo />
      </div>

      <div className="mb-6">
        <SmartData />
      </div>

      <div className="mb-6">
        <ContainerList />
      </div>
    </div>
  )
}

export default Dashboard

