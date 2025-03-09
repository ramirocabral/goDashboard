"use client";

const TimeRangeSelector = ({ timeRange, setTimeRange }) => {
  return (
    <div className="flex space-x-2 mt-2 sm:mt-0 overflow-x-auto w-full sm:w-auto justify-center">
      {["1h", "3h", "6h", "12h"].map((range) => (
        <button
          key={range}
          className={`px-3 py-1 text-sm font-medium rounded border whitespace-nowrap ${
            timeRange === range
              ? "bg-blue-950 text-white"
              : "bg-gray-700 text-gray-300 hover:bg-gray-600"
          }`}
          onClick={() => setTimeRange(range)}
        >
          {range}
        </button>
      ))}
    </div>
  );
};

export default TimeRangeSelector;