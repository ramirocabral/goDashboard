"use client";

import { useState } from "react";
import { useWebSocket } from "../contexts/WebSocketContext";
import CardContainer from "../components/widgets/WidgetContainer";
import { Box, Search, ChevronDown, ChevronUp } from "lucide-react";
const Containers = () => {
  const { containerData } = useWebSocket();
  const [filter, setFilter] = useState("");
  const [expanded, setExpanded] = useState(false);

  const filteredContainers =
    containerData?.filter(
      (container) =>
        container.name.toLowerCase().includes(filter.toLowerCase()) ||
        container.image.toLowerCase().includes(filter.toLowerCase()) ||
        container.id.toLowerCase().includes(filter.toLowerCase()),
    ) || [];

  const getTimeElapsed = (timestamp) => {
    const created = new Date(timestamp * 1000);
    const now = new Date();
    const diffMs = now.getTime() - created.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    const diffHours = Math.floor(
      (diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60),
    );
    if (diffDays > 0) return `${diffDays}d ${diffHours}h ago`;
    const diffMinutes = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60));
    return `${diffHours}h ${diffMinutes}m ago`;
  };

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case "running":
        return "text-green-500";
      case "paused":
        return "text-yellow-500";
      case "stopped":
        return "text-red-500";
      default:
        return "text-muted-foreground";
    }
  };

  return (
    <CardContainer>
      <div className="mb-4 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <div className="rounded-lg bg-cyan-500/10 p-2">
            <Box className="h-5 w-5 text-cyan-500" />
          </div>
          <div>
            <h3 className="text-sm font-medium text-gray-200">Containers</h3>
            <p className="text-xs text-gray-400">
              {containerData?.length || 0} containers running
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-2">
          <div className="relative">
            <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
            <input
              type="text"
              placeholder="Filter..."
              className="w-32 rounded border border-gray-700 bg-gray-800 py-1 pl-8 pr-2 text-xs text-gray-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
            />
          </div>
          <button
            className="rounded bg-gray-800 p-1 hover:bg-gray-700"
            onClick={() => setExpanded(!expanded)}
          >
            {expanded ? (
              <ChevronUp className="h-4 w-4 text-gray-400" />
            ) : (
              <ChevronDown className="h-4 w-4 text-gray-400" />
            )}
          </button>
        </div>
      </div>

      {/* loading screen  */}
      {!containerData ? (
        <div className="flex h-40 items-center justify-center">
          <div className="text-center text-gray-400">
            <Box className="mx-auto mb-2 h-8 w-8 animate-pulse" />
            <p>Loading container data...</p>
          </div>
        </div>
      ) : (
        <div className="overflow-x-auto">
          {filteredContainers.length === 0 ? (
            <div className="py-4 text-center text-gray-400">
              No containers found
            </div>
          ) : (
            <div className="space-y-2">
              <table className="w-full text-sm text-left text-gray-400">
                <thead className="bg-gray-700 text-gray-300">
                  <tr>
                    <th className="p-2">Name</th>
                    <th className="p-2">ID</th>
                    <th className="p-2">Status</th>
                    <th className="p-2">Image</th>
                    <th className="p-2">Created</th>
                    <th className="p-2">Uptime</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredContainers
                    .slice(0, expanded ? filteredContainers.length : 3)
                    .map((container) => (
                      <tr
                        key={container.id}
                        className="border-b border-gray-600 bg-gray-800/50"
                      >
                        <td className="p-2 font-medium text-gray-200">
                          {container.name}
                        </td>
                        <td className="p-2 font-mono">
                          {container.id.substring(0, 12)}
                        </td>
                        <td
                          className={`p-2 font-medium ${getStatusColor(container.status)}`}
                        >
                          {container.status}
                        </td>
                        <td className="p-2 text-gray-300">{container.image}</td>
                        <td className="p-2 text-gray-300">
                          {getTimeElapsed(container.created)}
                        </td>
                        <td className="p-2 text-gray-300">
                          {container.uptime}
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
              {filteredContainers.length > 3 && !expanded && (
                <button
                  className="mt-2 w-full rounded-md bg-gray-800 py-1 text-xs text-gray-400 hover:bg-gray-700"
                  onClick={() => setExpanded(true)}
                >
                  Show {filteredContainers.length - 3} more containers
                </button>
              )}

              {expanded && filteredContainers.length > 3 && (
                <button
                  className="mt-2 w-full rounded-md bg-gray-800 py-1 text-xs text-gray-400 hover:bg-gray-700"
                  onClick={() => setExpanded(false)}
                >
                  Show less
                </button>
              )}
            </div>
          )}
        </div>
      )}
    </CardContainer>
  );
};

export default Containers;
