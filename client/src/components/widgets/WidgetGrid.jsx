"use client"

const WidgetGrid = ({ data }) => (
    <div className="mb-4 grid grid-cols-2 gap-4">
      {data.map(({ label, value }, index) => (
        <div key={index}>
          <p className="text-xs text-gray-400">{label}</p>
          <p className="text-sm font-medium text-gray-200">{value}</p>
        </div>
      ))}
    </div>
);

export default WidgetGrid;