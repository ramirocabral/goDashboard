import { useWebSocket } from "../contexts/WebSocketContext"

const ContainerList = () => {
  const { containerData } = useWebSocket()

  if (!containerData || containerData.length === 0) {
    return (
      <div className="card">
        <div className="card-header">
          <div className="card-title">Containers</div>
          <div className="card-description">No containers running</div>
        </div>
        <div className="card-content">
          <div className="text-center text-muted-foreground py-4">No container data available</div>
        </div>
      </div>
    )
  }

  return (
    <div className="card">
      <div className="card-header">
        <div className="card-title">Containers</div>
        <div className="card-description">{containerData.length} containers running</div>
      </div>
      <div className="card-content">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border">
                <th className="text-left pb-2">Name</th>
                <th className="text-left pb-2">Status</th>
                <th className="text-left pb-2">Uptime</th>
                <th className="text-left pb-2">Image</th>
              </tr>
            </thead>
            <tbody>
              {containerData.map((container, index) => (
                <tr key={index} className="border-b border-border">
                  <td className="py-2">{container.name}</td>
                  <td className="py-2">
                    <span
                      className={`inline-block px-2 py-1 rounded-full text-xs ${
                        container.status.toLowerCase().includes("up")
                          ? "bg-green-500/20 text-green-500"
                          : "bg-yellow-500/20 text-yellow-500"
                      }`}
                    >
                      {container.status}
                    </span>
                  </td>
                  <td className="py-2">{container.uptime}</td>
                  <td className="py-2">{container.image}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

export default ContainerList

