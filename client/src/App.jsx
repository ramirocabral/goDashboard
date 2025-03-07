"use client"

import { useState, useEffect } from "react"
import { WebSocketProvider } from "./contexts/WebSocketContext"
import Dashboard from "./Dashboard"

function App() {
  const [darkMode, setDarkMode] = useState(false)

  // Apply dark mode class to body
  useEffect(() => {
    if (darkMode) {
      document.documentElement.classList.add("dark")
    } else {
      document.documentElement.classList.remove("dark")
    }
  }, [darkMode])

  return (
    <WebSocketProvider>
      <div className="min-h-screen bg-background text-foreground transition-colors duration-300">
        <Dashboard darkMode={darkMode} setDarkMode={setDarkMode} />
      </div>
    </WebSocketProvider>
  )
}

export default App

