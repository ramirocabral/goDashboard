"use client"

import React from "react";
import { ResponsiveContainer, AreaChart, Area } from "recharts";


const Chart = ({ realTimeData, id, stopColor, dataKey, stroke, strokeWidth, fill}) => (
    <div className="h-24">
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart data={realTimeData}>
          <defs>
            <linearGradient id={id} x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor={stopColor} stopOpacity={0.3} />
              <stop offset="100%" stopColor={stopColor} stopOpacity={0} />
            </linearGradient>
          </defs>
          <Area
            type="monotone"
            dataKey= {dataKey}
            stroke= {stroke}
            strokeWidth={strokeWidth}
            fill={fill}
            isAnimationActive={false}
            dot={false}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
    );

export default Chart;