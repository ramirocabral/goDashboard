"use client"

import Navbar from "../components/navbar/Navbar"

import HistoricalMetrics from "@/components/historical/HistoricalMetrics"

export default function History() {

  return (
    <div className="min-h-screen bg-background text-foreground transition-colors duration-300">
      <div className="container mx-auto p-4 sm:px-0">
      <Navbar />
      <HistoricalMetrics timeRange="1h"/>
      </div>
    </div>
  )
}