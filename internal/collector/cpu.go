package collector

import (
    "log"
    "strings"

    "golang-system-monitor/internal/utils"
)

type CPU struct{
    ModelName string            `json:"model_name"`
    Cores uint64                `json:"cores"`
    Threads uint64              `json:"threads"`
    Frequency float64           `json:"frequency"`
    Temp uint64                 `json:"temp"`
    UsageStatistics Usage       `json:"usage"`
}

type Usage struct{
    UsagePercentage float64     `json:"usage_percentage"`
    IdlePercentage float64      `json:"idle_percentage"`
    UserPercentage float64      `json:"user_percentage"`
    SystemPercentage float64    `json:"system_percentage"`  
    IoWaitPercentage float64    `json:"io_wait_percentage"`
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

    usage := Usage{}

    usage.UsagePercentage = cpuStats.UsageStatistics.IdlePercentage
    usage.IdlePercentage = cpuStats.UsageStatistics.IdlePercentage
    usage.UserPercentage = cpuStats.UsageStatistics.UserPercentage
    usage.SystemPercentage = cpuStats.UsageStatistics.SystemPercentage
    usage.IoWaitPercentage = cpuStats.UsageStatistics.IoWaitPercentage

    output.UsageStatistics = usage

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

    usage := Usage{}

    usage.UsagePercentage = (float64(total - idle64) / float64(total)) * 100
    usage.IdlePercentage = float64(idle64 / total) * 100
    usage.UserPercentage = float64(usr64 / total) * 100
    usage.SystemPercentage = float64(system64 / total) * 100
    usage.IoWaitPercentage = float64(iowait64 / total) * 100

    output.UsageStatistics = usage

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
