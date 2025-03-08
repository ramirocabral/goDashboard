"use client"

import { WebSocketProvider } from "./contexts/WebSocketContext"
import Dashboard from "./Dashboard"

function App() {
  return (
    <WebSocketProvider>
      <div className="min-h-screen bg-background text-foreground transition-colors duration-300">
        <Dashboard/>
      </div>
    </WebSocketProvider>
  )
}

export default App
