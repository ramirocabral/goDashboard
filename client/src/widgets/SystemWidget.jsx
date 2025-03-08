import React from 'react';
import { Github } from 'lucide-react';
import { useWebSocket } from '../contexts/WebSocketContext';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    faApple,
    faCentos,
    faDebian,
    faFedora,
    faLinux,
    faRedhat,
    faSuse,
    faUbuntu,
    faWindows,
  } from '@fortawesome/free-brands-svg-icons';

const SystemWidget = () => {
  const { systemInfo } = useWebSocket()
  const { uptimeData } = useWebSocket()

  const sysInfo = systemInfo

  return (
    <div className="shadow-sm border border-border bg-gradient-to-br from-gray-900 to-gray-800 rounded-xl p-8 h-full">
      <div className="flex flex-col">
        <div className="flex justify-between items-center">
          <a href="https://www.github.com/ramirocabral/golang-system-monitor" target='_blank' className="text-white text-sm flex items-center">
            <Github className="h-5 w-5 mr-2" />
          </a>
        </div>

      {!sysInfo ? (
        <div className="text-center text-gray-400 mt-4">Loading system data...</div>
      ):(

      <div className="relative flex flex-col items-center justify-center p-10 overflow-visible">
        <div className="absolute w-64 h-64 bg-gray-800/40 rounded-full flex items-center justify-center">
          <FontAwesomeIcon icon={getOsLogo( sysInfo.os.toLowerCase() )} className="w-48 h-48 opacity-15" />
        </div>
         

        <h1 className="text-white text-4xl font-bold z-10">{sysInfo.hostname}</h1>

        {/* System details */}
        <div className="mt-8 text-gray-300 text-sm z-10">
          <div className="grid grid-cols-2 gap-x-4 gap-y-1">
            <div className="text-gray-400">OS</div>
            <div>{sysInfo.os}</div>
            <div className="text-gray-400">Kernel</div>
            <div>{sysInfo.kernel}</div>
            <div className="text-gray-400">Uptime</div>
            {uptimeData ? (
              <div>{formatUptime(uptimeData.uptime)}</div>
            ) : (
              <div className="animate-pulse">Loading...</div>
            )}
          </div>
        </div>
      </div>
    )}
  </div>
</div>
  );
};

function getOsLogo(os){
  let icon = null;
  if (os.includes('ubuntu')) {
    icon = faUbuntu;
  } else if (os.includes('debian')){
    icon = faDebian
  } else if (os.includes('suse')) {
    icon = faSuse;
  } else if (os.includes('redhat')) {
    icon = faRedhat;
  } else if (os.includes('fedora')) {
    icon = faFedora;
  } else if (os.includes('centos')) {
    icon = faCentos;
  } else if (
    os.includes('mac') ||
    os.includes('osx') ||
    os.includes('darwin') ||
    os.includes('apple')
  ) {
    icon = faApple;
  } else if (os.includes('win')) {
    icon = faWindows;
  } else{
    icon = faLinux;
  }

  return icon;
}
const formatUptime = (seconds) => {
  if (!seconds) return "0 minutes"

  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)

  let result = ""
  if (days > 0) result += `${days} day${days > 1 ? "s" : ""} `
  if (hours > 0) result += `${hours} hour${hours > 1 ? "s" : ""} `
  if (minutes > 0) result += `${minutes} minute${minutes > 1 ? "s" : ""}`

  return result.trim()
}


export default SystemWidget;