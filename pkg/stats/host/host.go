package host

import (
    "errors"
    "log"
    "strings"
    "time"

    "golang-system-monitor/internal/utils"
)

type Host struct{
    Hostname    string  `json:"hostname"`
    Os          string  `json:"os"`
    Kernel      string  `json:"kernel"`
    Date        string  `json:"date"`
}

const HOSTNAME_PATH = "/proc/sys/kernel/hostname"
const OS_PATH = "/etc/os-release"
const KERNEL_PATH = "/proc/version"
const LAST_BOOT_PATH = "/proc/stat"

func ReadHost() (Host, error){
    output := Host{}

    // get hostname
    hostnameData, err := utils.ReadFile(HOSTNAME_PATH)
    if err != nil{
        log.Println("Error reading hostname data: ", err)
        return output, errors.New("error reading hostname data")

    }
    output.Hostname = strings.TrimSpace(string(hostnameData))

    // get os release
    osData, err := utils.ReadFile(OS_PATH)

    if err != nil{
        log.Println("Error reading os data: ", err)
        return output, errors.New("error reading os data")
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
        return output, errors.New("error reading kernel data")
    }

    kernelDataSplit := strings.Split(string(kernelData), " ")
    output.Kernel = kernelDataSplit[2]


    //get date

    date := time.Now().UTC().Format(time.RFC3339)

    output.Date = date


    return output, nil
}
