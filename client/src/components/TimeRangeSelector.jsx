const TimeRangeSelector = ({ selectedRange, onRangeChange }) => {
  const ranges = [
    { label: "1h", value: 3600 },
    { label: "3h", value: 10800 },
    { label: "6h", value: 21600 },
    { label: "12h", value: 43200 },
  ]

  return (
    <div className="flex space-x-2">
      {ranges.map((range) => (
        <button
          key={range.value}
          className={`time-range-button ${selectedRange === range.value ? "active" : ""}`}
          onClick={() => onRangeChange(range.value)}
        >
          {range.label}
        </button>
      ))}
    </div>
  )
}

export default TimeRangeSelector

