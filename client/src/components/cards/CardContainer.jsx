"use client"

const CardContainer = ({ children }) => (
    <div className="relative overflow-hidden rounded-xl bg-gradient-to-br from-gray-900 to-gray-800 p-4 shadow-lg transition-all hover:shadow-xl">
      {children}
    </div>
);

export default CardContainer