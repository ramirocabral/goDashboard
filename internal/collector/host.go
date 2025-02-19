package stats

import (
	"errors"
	"log"
	"strings"

	"golang-system-monitor/internal/utils"
)


type Host struct{
    Hostname string
    Os string
    Kernel string
    LastBoot string
    Date string
}

const HOSTNAME_PATH = "/proc/sys/kernel/hostname"
const OS_PATH = "/etc/os-release"
const KERNEL_PATH = "/proc/version"
const LAST_BOOT_PATH = "/proc/stat"
const DATE_PATH = "/proc/driver/rtc"

func ReadHost() (Host, error){
    output := Host{}

    // get hostname
    hostnameData, err := utils.ReadFile(HOSTNAME_PATH)
    if err != nil{
        log.Println("Error reading hostname data: ", err)
        output.Hostname = "N/A"
    }
    output.Hostname = strings.TrimSpace(string(hostnameData))

    // get os release
    osData, err := utils.ReadFile(OS_PATH)

    if err != nil{
        log.Println("Error reading os data: ", err)
        output.Os = "N/A"
    }

    osDataSplit := strings.Split(string(osData), "\n")

    for _, line := range osDataSplit{
        fields := strings.Fields(line)
        if len(fields) == 0{
            continue
        }
        if strings.Contains(fields[0], "PRETTY_NAME"){
            output.Os = strings.Trim(fields[1], "\"")
        }
    }

    // get kernel version
    kernelData, err := utils.ReadFile(KERNEL_PATH)

    if err != nil{
        log.Println("Error reading kernel data: ", err)
        output.Kernel = "N/A"
    }else{
        kernelDataSplit := strings.Split(string(kernelData), " ")
        output.Kernel = kernelDataSplit[2]

    }
    // get date
    dateData, err := utils.ReadFile(DATE_PATH)
    
    if err != nil{
        log.Println("Error reading date data: ", err)
        output.Date = "N/A"
    }
    else{
        dateDataSplit := strings.Split(string(dateData), "\n")
    }








        


    return output, nil
}
