// "use client"

// import { Link, useLocation } from "react-router-dom";
// import { useWebSocket } from "@/contexts/WebSocketContext";
// import { APP_NAME } from "../../constants/config";

// const Navbar = () => {
//     const location = useLocation();

//     const { connected } = useWebSocket();
//     const isConnected = Object.values(connected).some((status) => status);
//     const isOnCharts = location.pathname === "/charts";

//     return (
//         <header className="flex justify-between items-center mb-5">
//         <div className="flex items-center">
//           <h1 className="text-3xl font-bold">{APP_NAME}</h1>
//           <div className="flex items-center ml-4"></div>
//           <div
//             className={`ml-4 flex items-center ${
//               isConnected ? "text-green-500" : "text-red-500"
//             }`}
//           >
//             <div
//               className={`h-2 w-2 rounded-full ${
//                 isConnected ? "bg-green-500 animate-pulse" : "bg-red-500"
//               }`}
//             />
//             <span className="ml-2 text-sm">
//               {isConnected ? "Connected" : "Disconnected"}
//             </span>
//           </div>
//         </div>
//         <div className="flex align-middle">
//           <Link
//             to={isOnCharts ? "/" : "/charts"}
//             className="btn-custom"
//           >
//             {isOnCharts ? "Back to Dashboard" : "Metrics Charts"}
//           </Link>
//         </div>
//       </header>
//     )
// }

// export default Navbar;

"use client";

import { Link, useLocation } from "react-router-dom";
import { useWebSocket } from "@/contexts/WebSocketContext";
import { APP_NAME } from "../../constants/config";

const Navbar = () => {
    const location = useLocation();
    const { connected } = useWebSocket();
    const isConnected = Object.values(connected).some((status) => status);
    const isOnCharts = location.pathname === "/charts";

    return (
        <header className="flex flex-wrap sm:flex-nowrap sm:justify-between justify-center items-center mb-5">
            <div className="flex items-center w-full sm:w-auto justify-between sm:justify-start">
                <h1 className="text-xl sm:text-3xl font-bold sm:block">{APP_NAME}</h1>

                <div className="ml-4 flex items-center text-sm sm:text-base">
                    <div className={`h-3 w-3 rounded-full ${isConnected ? "bg-green-500 animate-pulse" : "bg-red-500"}`} />
                    <span className="ml-2">{isConnected ? "Connected" : "Disconnected"}</span>
                </div>
            </div>

            <div className="mt-3 sm:mt-0 w-full sm:w-auto flex justify-center">
                <Link to={isOnCharts ? "/" : "/charts"} className="btn-custom px-4 py-2 text-sm sm:text-base">
                    {isOnCharts ? "Back to Dashboard" : "Metrics Charts"}
                </Link>
            </div>
        </header>
    );
};

export default Navbar;
