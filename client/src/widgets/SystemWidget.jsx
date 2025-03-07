import React from 'react';
import { Github } from 'lucide-react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    faApple,
    faCentos,
    faDebian,
    faFedora,
    faGithub,
    faLinux,
    faRedhat,
    faSuse,
    faUbuntu,
    faWindows,
  } from '@fortawesome/free-brands-svg-icons';

const SystemWidget = () => {
    
  const systemInfo = {
    name: "dash.",
    os: "Ubuntu 22.04 LTS",
    arch: "x64",
    uptime: "27 minutes"
  };

  const osLogo = getOsLogo("debian");

  // State for dark mode toggle
  const [darkMode, setDarkMode] = React.useState(true);

  return (
    <div className="shadow-sm border border-border bg-gradient-to-br from-gray-900 to-gray-800 rounded-xl p-8">
      <div className="flex flex-col">
        {/* Top section with GitHub icon */}
        <div className="flex justify-between items-center">
          <Github className="text-white h-5 w-5" />
        </div>

      {/* Center content with circular background */}
      <div className="relative flex flex-col items-center justify-center p-10 overflow-visible">
        <div className="absolute w-64 h-64 bg-gray-800/40 rounded-full flex items-center justify-center">
          <FontAwesomeIcon icon={osLogo} className="w-48 h-48 opacity-15" />
        </div>
         

        <h1 className="text-white text-4xl font-bold z-10">gelsomina</h1>

        {/* System details */}
        <div className="mt-8 text-gray-300 text-sm z-10">
          <div className="grid grid-cols-2 gap-x-4 gap-y-1">
            <div className="text-gray-400">OS</div>
            <div>ArchLinux</div>
            <div className="text-gray-400">Arch</div>
            <div>x64</div>
            <div className="text-gray-400">Up since</div>
            <div>2h</div>
          </div>
        </div>
      </div>
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
  } else if (os.includes('linux')) {
    icon = faLinux;
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

export default SystemWidget;