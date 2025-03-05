"use client"

import React from "react"

export const Tabs = ({ defaultValue, className, children }) => {
  const [activeTab, setActiveTab] = React.useState(defaultValue)

  // Clone children and pass activeTab to them
  const childrenWithProps = React.Children.map(children, (child) => {
    if (React.isValidElement(child)) {
      return React.cloneElement(child, { activeTab, setActiveTab })
    }
    return child
  })

  return <div className={`tabs ${className || ""}`}>{childrenWithProps}</div>
}

export const TabsList = ({ className, children, activeTab, setActiveTab }) => {
  return <div className={`flex space-x-2 border-b border-border ${className || ""}`}>{children}</div>
}

export const TabsTrigger = ({ value, children, activeTab, setActiveTab }) => {
  const isActive = activeTab === value

  return (
    <button
      className={`px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors ${
        isActive
          ? "border-primary text-primary"
          : "border-transparent text-muted-foreground hover:text-foreground hover:border-border"
      }`}
      onClick={() => setActiveTab(value)}
    >
      {children}
    </button>
  )
}

export const TabsContent = ({ value, children, activeTab }) => {
  if (activeTab !== value) return null

  return <div className="py-4">{children}</div>
}