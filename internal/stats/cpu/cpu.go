package cpu

import (
	"log"
	"strings"
	"time"

	"golang-system-monitor/internal/utils"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPU struct{
    ModelName string            `json:"model_name"`
    Cores uint64                `json:"cores"`
    Threads uint64              `json:"threads"`
    Temp uint64                 `json:"temp"`
    UsageStatistics Usage       `json:"usage"`
}

type Usage struct{
    UsagePercentage float64     `json:"usage_percentage"`
    IdlePercentage float64      `json:"idle_percentage"`
}

func (c CPU) ToMap() map[string]interface{}{
    return map[string]interface{}{
        "model_name": c.ModelName,
        "cores": c.Cores,
        "threads": c.Threads,
        "temp": c.Temp,
        "usage_percentage": c.UsageStatistics.UsagePercentage,
        "idle_percentage": c.UsageStatistics.IdlePercentage,
    }
}

var firstRun = true

const CPU_PATH = "/proc/stat"
const CPU_INFO_PATH = "/proc/cpuinfo"

func ReadCPU() (CPU, error){
    output := CPU{}

    //get percentages
    percentages, err := cpu.Percent(1 * time.Second, false)

    if err != nil{
        log.Println("Error getting CPU percentages: ", err)
        return CPU{}, err
    }

    usage := Usage{}

    //reformat without two decimal points
    usage.UsagePercentage = utils.RoundFloat(percentages[0], 2)
    usage.IdlePercentage = utils.RoundFloat(100 - percentages[0], 2)

    output.UsageStatistics = usage

    //get cpu temperature
    tmp,err := getCpuTemp()

    if err != nil{
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

    return output, nil
}

// func getCPUStats()(CPU, error){
//     output := CPU{}
//
//     cpuData, err := utils.ReadFile(CPU_PATH)
//     if err != nil{
//         log.Println("Error reading CPU data: ", err)
//         return output, err
//     }
//
//     cpuDataStr := strings.Split((cpuData), "\n")
//
//     core_data := strings.Fields(cpuDataStr[0])
//
//     usr, nice, system, idle, iowait, irq, softirq := core_data[1], core_data[2], core_data[3], core_data[4], core_data[5], core_data[6], core_data[7]
//
//     usr64 := utils.StrToUint64(usr)
//     nice64 := utils.StrToUint64(nice)
//     system64 := utils.StrToUint64(system)
//     idle64 := utils.StrToUint64(idle)
//     iowait64 := utils.StrToUint64(iowait)
//     irq64 := utils.StrToUint64(irq)
//     softirq64 := utils.StrToUint64(softirq)
//
//     total := usr64 + nice64 + system64 + idle64 + iowait64 + irq64 + softirq64
//
//     fmt.Println(total)
//
//     usage := Usage{}
//
//     usage.UsagePercentage = (float64(total - idle64) / float64(total)) * 100
//
//     fmt.Println(usage.UsagePercentage)
//     usage.IdlePercentage = (float64(idle64) / float64(total)) * 100
//     fmt.Println(usage.IdlePercentage)
//     usage.UserPercentage = float64(usr64 / total) * 100
//     usage.SystemPercentage = float64(system64 / total) * 100
//     usage.IoWaitPercentage = float64(iowait64 / total) * 100
//
//     output.UsageStatistics = usage
//
//     return output, nil
// }

func getCpuTemp() (uint64, error){
    tempFile, err := utils.FindCPUTempFile()
    

    //encuentra el archivos
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
    }

    return output, nil
}
