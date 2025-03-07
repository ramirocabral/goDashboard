package cpu

import (
	"log"
	"strings"
	"time"

	"golang-system-monitor/internal/core"
	"golang-system-monitor/internal/utils"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPU struct{
    ModelName       string      `json:"model_name"`
    Frequency       uint64      `json:"frequency"`
    Family          string      `json:"family"`
    Cores           uint64      `json:"cores"`
    Threads         uint64      `json:"threads"`
    Temp            uint64      `json:"temp"`
    UsageStatistics Usage       `json:"usage"`
}

type Usage struct{
    UsagePercentage float64     `json:"usage_percentage"`
    IdlePercentage float64      `json:"idle_percentage"`
}

func (c CPU) ToPoint() []*core.Point{
    return []*core.Point{{
        Timestamp: time.Now(),
        Measurement: "cpu",
        Tags: map[string]string{
            "model_name": c.ModelName,
            "frequency": utils.Uint64ToStr(c.Frequency),
            "family": c.Family,
            "cores": utils.Uint64ToStr(c.Cores),
            "threads": utils.Uint64ToStr(c.Threads),
        },
        Fields: map[string]interface{}{
            "temp" : c.Temp,
            "usage_percentage": c.UsageStatistics.UsagePercentage,
            "idle_percentage": c.UsageStatistics.IdlePercentage,
        },
    } ,
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

    output.Frequency = cpuInfo.Frequency
    output.Family = cpuInfo.Family
    output.ModelName = cpuInfo.ModelName
    output.Cores = cpuInfo.Cores
    output.Threads = cpuInfo.Threads

    return output, nil
}

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

    lines, err := utils.ExecuteCommand("dmidecode", "-t", "4")

    if err != nil{
        log.Println("Error reading CPU info: ", err)
        return CPU{}, err
    }

    for _, line := range strings.Split(lines, "\n"){
        if strings.Contains(line, "Max Speed"){
            output.Frequency = utils.StrToUint64(strings.Fields(line)[2])
        }
        if strings.HasPrefix(strings.TrimSpace(line), "Family"){
            output.Family = strings.TrimSpace(strings.Split(line, ":")[1])
        }
    }

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
