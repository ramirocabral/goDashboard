import React from "react"
import ReactDOM from "react-dom/client"
import App from "./App.jsx"
import "./index.css"
import { WebSocketProvider } from "./contexts/WebSocketContext.jsx"

ReactDOM.createRoot(document.getElementById("root")).render(
  <WebSocketProvider>
    <App />
  </WebSocketProvider>
)