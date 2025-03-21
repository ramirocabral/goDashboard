"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";
import {
  API_BASE_URL_REALTIME,
  API_BASE_URL_STATS,
} from "../constants/endpoints";

const WebSocketContext = createContext(null);

// WebSocket endpoints
const WS_ENDPOINTS = {
  CPU: `${API_BASE_URL_REALTIME}/cpu`,
  MEMORY: `${API_BASE_URL_REALTIME}/memory`,
  NETWORK: `${API_BASE_URL_REALTIME}/network`,
  CONTAINER: `${API_BASE_URL_REALTIME}/container`,
  UPTIME: `${API_BASE_URL_REALTIME}/uptime`,
  IO: `${API_BASE_URL_REALTIME}/io`,
};

// REST API endpoints
const API_ENDPOINTS = {
  SYSTEM_INFO: `${API_BASE_URL_STATS}/host`,
  DISKS: `${API_BASE_URL_STATS}/disk`,
  SMART: `${API_BASE_URL_STATS}/smart`,
};

export const WebSocketProvider = ({ children }) => {
  const [cpuData, setCpuData] = useState(null);
  const [memoryData, setMemoryData] = useState(null);
  const [networkData, setNetworkData] = useState(null);
  const [containerData, setContainerData] = useState(null);
  const [uptimeData, setUptimeData] = useState(null);
  const [ioData, setIoData] = useState(null);

  const [systemInfo, setSystemInfo] = useState(null);
  const [disksInfo, setDisksInfo] = useState(null);
  const [smartData, setSmartData] = useState(null);

  const [connected, setConnected] = useState({
    cpu: false,
    memory: false,
    network: false,
    container: false,
    uptime: false,
    io: false,
  });

  const createWebSocket = useCallback((url, dataType) => {
    const socket = new WebSocket(url);

    socket.onopen = () => {
      console.log(`${dataType} WebSocket connected`);
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: true }));
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        switch (data.type) {
          case "cpu":
            setCpuData(data.data);
            break;
          case "memory":
            setMemoryData(data.data);
            break;
          case "network":
            setNetworkData(data.data);
            break;
          case "container":
            setContainerData(data.data);
            break;
          case "uptime":
            setUptimeData(data.data);
            break;
          case "io":
            setIoData(data.data);
            break;
          default:
            console.warn("Unknown data type:", data.type);
        }
      } catch (error) {
        console.error("Error parsing WebSocket data:", error);
      }
    };

    socket.onclose = () => {
      console.log(`${dataType} WebSocket disconnected`);
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: false }));

      setTimeout(() => {
        createWebSocket(url, dataType);
      }, 5000);
    };

    socket.onerror = (error) => {
      console.error(`${dataType} WebSocket error:`, error);
      setConnected((prev) => ({ ...prev, [dataType.toLowerCase()]: false }));
    };

    return socket;
  }, []);

  // create WebSocket connections on component mount
  useEffect(() => {
    const sockets = {
      cpu: createWebSocket(WS_ENDPOINTS.CPU, "CPU"),
      memory: createWebSocket(WS_ENDPOINTS.MEMORY, "MEMORY"),
      network: createWebSocket(WS_ENDPOINTS.NETWORK, "NETWORK"),
      container: createWebSocket(WS_ENDPOINTS.CONTAINER, "CONTAINER"),
      uptime: createWebSocket(WS_ENDPOINTS.UPTIME, "UPTIME"),
      io: createWebSocket(WS_ENDPOINTS.IO, "IO"),
    };

    // clean up WebSocket connections on unmount
    return () => {
      Object.values(sockets).forEach((socket) => {
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.close();
        }
      });
    };
  }, [createWebSocket]);

  // fetch system info, disks info, and SMART data
  useEffect(() => {
    const fetchSystemInfo = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.SYSTEM_INFO);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setSystemInfo(data);
      } catch (error) {
        console.error("Error fetching system info:", error);
      }
    };

    const fetchDisksInfo = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.DISKS);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setDisksInfo(data);
      } catch (error) {
        console.error("Error fetching disks info:", error);
      }
    };

    const fetchSmartData = async () => {
      try {
        const response = await fetch(API_ENDPOINTS.SMART);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        setSmartData(data);
      } catch (error) {
        console.error("Error fetching SMART data:", error);
      }
    };

    fetchSystemInfo();
    fetchDisksInfo();
    fetchSmartData();

    const interval = setInterval(
      () => {
        fetchSystemInfo();
        fetchDisksInfo();
        fetchSmartData();
      },
      5 * 60 * 1000,
    );

    return () => clearInterval(interval);
  }, []);

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
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};
