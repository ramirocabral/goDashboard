"use client"

import Dashboard from "./pages/Dashboard"
import Historical from "./pages/Charts"
import { Route, BrowserRouter as Router, Routes } from "react-router-dom"

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/charts" element={<Historical />} />
      </Routes>
    </Router>
  )
}

export default App
