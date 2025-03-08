"use client";

import { useState } from "react";
import { ChevronDown, ChevronUp } from "lucide-react";

const WidgetHeaderSelector = ({
  icon,
  title,
  items = [],
  selectedItem,
  onItemSelect,
  valueText,
  valueSubtext,
  selectorLabel = null,
}) => {
  const [showSelector, setShowSelector] = useState(false);

  const showDropdown = items.length > 1;

  return (
    <div className="mb-4 flex items-center justify-between">
      <div className="flex items-center space-x-3">
        {icon}
        <div>
          <h3 className="text-sm font-medium text-gray-200">{title}</h3>
          <div className="relative">
            <button
              onClick={() => setShowSelector(!showSelector)}
              className="flex items-center text-xs text-gray-400 hover:text-gray-300"
            >
              {selectorLabel || selectedItem}
              {showDropdown &&
                (showSelector ? (
                  <ChevronUp className="ml-1 h-3 w-3" />
                ) : (
                  <ChevronDown className="ml-1 h-3 w-3" />
                ))}
            </button>

            {showSelector && showDropdown && (
              <div className="absolute top-full left-0 z-10 mt-1 w-40 rounded-md bg-gray-800 shadow-lg">
                <ul className="py-1">
                  {items.map((item) => (
                    <li key={item}>
                      <button
                        className={`block w-full px-4 py-2 text-left text-xs ${
                          item === selectedItem
                            ? "bg-gray-700 text-gray-200"
                            : "text-gray-400 hover:bg-gray-700 hover:text-gray-200"
                        }`}
                        onClick={() => {
                          onItemSelect(item);
                          setShowSelector(false);
                        }}
                      >
                        {item}
                      </button>
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </div>
      </div>
      {valueText && (
        <div className="text-right">
          <p className="text-2xl font-bold text-gray-200">{valueText}</p>
          {valueSubtext && (
            <p className="text-xs text-gray-400">{valueSubtext}</p>
          )}
        </div>
      )}
    </div>
  );
};

export default WidgetHeaderSelector;
