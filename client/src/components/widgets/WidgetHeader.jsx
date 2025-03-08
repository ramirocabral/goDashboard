"use client"

const CardHeader = ({ icon, title, subtitle, value, secondValue }) => (
  <div className="mb-4 flex items-center justify-between">
    <div className="flex items-center space-x-3">
      <div className="rounded-lg bg-blue-500/10 p-2">{icon}</div>
      <div>
        <h3 className="text-sm font-medium text-gray-200">{title}</h3>
        <p className="text-xs text-gray-400">{subtitle}</p>
      </div>
    </div>
    <div className="text-right">
      <p className="text-2xl font-bold text-gray-200">{value}</p>
      <p className="text-xs text-gray-400">{secondValue}</p>
    </div>
  </div>
);

export default CardHeader