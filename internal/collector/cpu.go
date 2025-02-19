package collector

import (
    "log"
    "strings"

    "golang-system-monitor/internal/utils"
)

type CPU struct{
    ModelName string
    Cores uint64
    Threads uint64
    Frequency float64
    Temp uint64                     //CPU temperature
    UsagePercentage float64         //Percentage of CPU used
    IdlePercentage float64          //Percentage of time spent idle
    UserPercentage float64          //Percentage of time spent running user code
    SystemPercentage float64        //Percentage of time spent running system code
    IoWaitPercentage float64        //Percentage of time spent waiting for I/O operations to complete
}

const CPU_PATH = "/proc/stat"
const CPU_INFO_PATH = "/proc/cpuinfo"

func ReadCPU() (CPU, error){
    output := CPU{}

    //get cpu stats
    cpuStats, err := getCPUStats()

    if err != nil{
        log.Println("Error getting CPU stats: ", err)
        return CPU{}, err
    }

    output.UsagePercentage = cpuStats.UsagePercentage
    output.IdlePercentage = cpuStats.IdlePercentage
    output.UserPercentage = cpuStats.UserPercentage
    output.SystemPercentage = cpuStats.SystemPercentage
    output.IoWaitPercentage = cpuStats.IoWaitPercentage

    //get cpu temperature
    tmp,err := getCpuTemp()

    if err == nil{
        log.Println("Error getting CPU temp: ", err)
        return CPU{}, err
    }
    output.Temp = tmp


    //get cpu model info
    cpuInfo, err := getCpuInfo()
    if err != nil{
        log.Println("Error getting CPU info: ", err)
        return CPU{}, err
    }

    output.ModelName = cpuInfo.ModelName
    output.Cores = cpuInfo.Cores
    output.Threads = cpuInfo.Threads
    output.Frequency = cpuInfo.Frequency

    return output, nil
}

func getCPUStats()(CPU, error){
    output := CPU{}

    cpuData, err := utils.ReadFile(CPU_PATH)
    if err != nil{
        log.Println("Error reading CPU data: ", err)
        return output, err
    }

    cpuDataStr := strings.Split((cpuData), "\n")

    core_data := strings.Fields(cpuDataStr[0])

    usr, nice, system, idle, iowait, irq, softirq := core_data[1], core_data[2], core_data[3], core_data[4], core_data[5], core_data[6], core_data[7]

    usr64 := utils.StrToUint64(usr)
    nice64 := utils.StrToUint64(nice)
    system64 := utils.StrToUint64(system)
    idle64 := utils.StrToUint64(idle)
    iowait64 := utils.StrToUint64(iowait)
    irq64 := utils.StrToUint64(irq)
    softirq64 := utils.StrToUint64(softirq)

    total := usr64 + nice64 + system64 + idle64 + iowait64 + irq64 + softirq64

    output.UsagePercentage = (float64(total - idle64) / float64(total)) * 100
    output.IdlePercentage = float64(idle64 / total) * 100
    output.UserPercentage = float64(usr64 / total) * 100
    output.SystemPercentage = float64(system64 / total) * 100
    output.IoWaitPercentage = float64(iowait64 / total) * 100

    return output, nil
}

func getCpuTemp() (uint64, error){
    tempFile, err := utils.FindCPUTempFile()

    if err != nil{
        log.Println("Error finding CPU temp file: ", err)
        return 0, err
    }

    temp, err := utils.ReadFile(tempFile)
    if err != nil{
        log.Println("Error reading temp file: ", err)
        return 0, err
    }

    tempStr := strings.TrimSpace(temp)

    cpuTemp := utils.StrToUint64(tempStr[:2])

    return cpuTemp, nil
}

func getCpuInfo() (CPU, error){
    output := CPU{}

    cpuInfo, err := utils.ReadFile(CPU_INFO_PATH)

    if err != nil{
        log.Println("Error reading CPU info: ", err)
        return output, err
    }

    cpuInfoStr := strings.Split(cpuInfo, "\n")

    for _, line := range cpuInfoStr{
        if strings.Contains(line, "model name"){
            output.ModelName = strings.TrimSpace(strings.Split(line, ":")[1])
        }
        if strings.Contains(line, "cpu cores"){
            output.Cores = utils.StrToUint64(strings.TrimSpace(strings.Split(line, ":")[1]))
        }
        if strings.Contains(line, "siblings"){
            output.Threads = utils.StrToUint64(strings.TrimSpace(strings.Split(line, ":")[1]))
        }
        if strings.Contains(line, "cpu MHz"){
            output.Frequency = utils.StrToFloat64(strings.TrimSpace(strings.Split(line, ":")[1]))
        }
    }

    return output, nil
}
