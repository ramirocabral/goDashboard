"use client"

import { createContext, useContext, useState, useEffect, useCallback } from "react"

// Create context
const WebSocketContext = createContext(null)

const BASE_URL="http://localhost:8080/api/v1"

// WebSocket endpoints
const WS_ENDPOINTS = {
  CPU: `${BASE_URL}/ws/cpu`,
  MEMORY: `${BASE_URL}/ws/memory`,
  NETWORK: `${BASE_URL}/ws/network`,
  CONTAINER: `${BASE_URL}/ws/container`,
  UPTIME: `${BASE_URL}/ws/uptime`,
  IO: `${BASE_URL}/ws/io`,
}

// REST API endpoints
const API_ENDPOINTS = {
  SYSTEM_INFO: `${BASE_URL}/stat/host`,
  DISKS: `${BASE_URL}/stat/disk`,
  SMART: `${BASE_URL}/stat/smart`,
}

export const WebSocketProvider = ({ children }) => {
  // State for each type of data
  const [cpuData, setCpuData] = useState(null)
  const [memoryData, setMemoryData] = useState(null)
  const [networkData, setNetworkData] = useState(null)
  const [containerData, setContainerData] = useState(null)
  const [uptimeData, setUptimeData] = useState(null)
  const [ioData, setIoData] = useState(null)

  // State for REST API data
  const [systemInfo, setSystemInfo] = useState(null)
  const [disksInfo, setDisksInfo] = useState(null)
  const [smartData, setSmartData] = useState(null)

  // Connection status
  const [connected, setConnected] = useState({
    cpu: false,
    memory: false,
    network: false,
    container: false,
    uptime: false,
    io: false,
  })

  // Function to create WebSocket connections
  const createWebSocket = useCallback((url, dataType) => {
    const socket = new WebSocket(url)

    socket.onopen = () => {
      console.log(`${dataType} WebSocket connected`)
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: true }))
    }

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)

        switch (data.type) {
          case "cpu":
            setCpuData(data.data)
            break
          case "memory":
            setMemoryData(data.data)
            break
          case "network":
            setNetworkData(data.data)
            break
          case "container":
            setContainerData(data.data)
            break
          case "uptime":
            setUptimeData(data.data)
            break
          case "io":
            setIoData(data.data)
            break
          default:
            console.warn("Unknown data type:", data.type)
        }
      } catch (error) {
        console.error("Error parsing WebSocket data:", error)
      }
    }

    socket.onclose = () => {
      console.log(`${dataType} WebSocket disconnected`)
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: false }))

      // Attempt to reconnect after 5 seconds
      setTimeout(() => {
        createWebSocket(url, dataType)
      }, 5000)
    }

    socket.onerror = (error) => {
      console.error(`${dataType} WebSocket error:`, error)
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: false }))
    }

    return socket
  }, [])

  // create WebSocket connections on component mount
  useEffect(() => {
    const sockets = {
      cpu: createWebSocket(WS_ENDPOINTS.CPU, "CPU"),
      memory: createWebSocket(WS_ENDPOINTS.MEMORY, "MEMORY"),
      network: createWebSocket(WS_ENDPOINTS.NETWORK, "NETWORK"),
      container: createWebSocket(WS_ENDPOINTS.CONTAINER, "CONTAINER"),
      uptime: createWebSocket(WS_ENDPOINTS.UPTIME, "UPTIME"),
      io: createWebSocket(WS_ENDPOINTS.IO, "IO"),
    }

    // clean up WebSocket connections on unmount
    return () => {
      Object.values(sockets).forEach((socket) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.close()
        }
      })
    }
  }, [createWebSocket])

  // fetch system info, disks info, and SMART data
  useEffect(() => {
    const fetchSystemInfo = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.SYSTEM_INFO)
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        const data = await response.json()
        setSystemInfo(data)
      } catch (error) {
        console.error("Error fetching system info:", error)
      }
    }

    const fetchDisksInfo = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.DISKS)
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        const data = await response.json()
        setDisksInfo(data)
      } catch (error) {
        console.error("Error fetching disks info:", error)
      }
    }

    const fetchSmartData = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.SMART)
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        const data = await response.json()
        // console.log(smartData.devices[0].device)
        setSmartData(data)
      } catch (error) {
        console.error("Error fetching SMART data:", error)
      }
    }

    fetchSystemInfo()
    fetchDisksInfo()
    fetchSmartData()

    // refresh the data every 5 minutes
    const interval = setInterval(
      () => {
        fetchSystemInfo()
        fetchDisksInfo()
        fetchSmartData()
      },
      5 * 60 * 1000,
    )

    return () => clearInterval(interval)
  }, [])

  // Function to fetch historical data
  const fetchHistoricalData = async (dataType, startTime, endTime, interval) => {
    try {
      const response = await fetch(
        `${BASE_URL}/api/${dataType}/history?start=${startTime}&end=${endTime}&interval=${interval}`,
      )
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      const data = await response.json()
      return data
    } catch (error) {
      console.error(`Error fetching ${dataType} historical data:`, error)
      return null
    }
  }

  return (
    <WebSocketContext.Provider
      value={{
        cpuData,
        memoryData,
        networkData,
        containerData,
        uptimeData,
        ioData,
        systemInfo,
        disksInfo,
        smartData,
        connected,
        fetchHistoricalData,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  )
}

// Custom hook to use the WebSocket context
export const useWebSocket = () => {
  const context = useContext(WebSocketContext)
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider")
  }
  return context
}